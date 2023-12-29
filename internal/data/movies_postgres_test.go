package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host   = "localhost"
	user   = "postgres"
	pw     = "postgres"
	dbName = "movies_test"
	port   = "5435"
	dsn    = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo MovieModel

func TestMain(m *testing.M) {
	// connect to docker; fail if docker not running
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	pool = p

	// set up our docker options, specifying the image and so forth
	options := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + pw,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get a resource (docker image)
	resource, err = pool.RunWithOptions(&options)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until it's ready
	if err = pool.Retry(func() error {
		testDB, err = sql.Open("postgres", fmt.Sprintf(dsn, host, port, user, pw, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	// populate the database with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepo = MovieModel{DB: testDB}

	// run tests
	code := m.Run()

	// clean up
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/movies.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func TestPingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

var movieID uuid.UUID

func TestPostgresDBRepoInsertMovie(t *testing.T) {
	testMovie := &Movie{
		Title:   "The Lord of the Rings: The Fellowship of the Ring",
		Year:    2001,
		Runtime: 178,
		Genres:  []string{"Action", "Adventure", "Drama", "Fantasy"},
	}

	err := testRepo.Insert(testMovie)
	if err != nil {
		t.Errorf("insert movie returned error: %s", err)
	}

	movieID = testMovie.ID
}

func TestPostgresDBRepoGetMovie(t *testing.T) {
	movie, err := testRepo.Get(movieID)
	if err != nil {
		t.Errorf("cannot get movie: %s", err)
	}

	if movie.Title != "The Lord of the Rings: The Fellowship of the Ring" {
		t.Errorf("expected 'The Lord of the Rings: The Fellowship of the Ring' but got %s", movie.Title)
	}

	if len(movie.Genres) != 4 {
		t.Errorf("expected 4 genres got but %d", len(movie.Genres))
	}
}

func TestPostgresDBRepoUpdateMovie(t *testing.T) {
	movie := &Movie{
		ID:      movieID,
		Title:   "The Lord of the Rings: The Two Towers",
		Genres:  []string{"Action", "Adventure", "Drama", "Fantasy", "Mystery"},
		Year:    2002,
		Version: 1,
	}

	err := testRepo.Update(movie)
	if err != nil {
		t.Errorf("cannot update movie: %s", err)
	}

	updatedMovie, err := testRepo.Get(movieID)
	if err != nil {
		t.Errorf("cannot get movie: %s", err)
	}

	if movie.Title != updatedMovie.Title {
		t.Errorf("expected %s but got %s", movie.Title, updatedMovie.Title)
	}

	if updatedMovie.Version != 2 {
		t.Errorf("expected version %d", updatedMovie.Version)
	}

	if len(movie.Genres) != 5 {
		t.Errorf("expected 5 genres got but %d", len(movie.Genres))
	}
}

func TestPostgresDBRepoGetAll(t *testing.T) {
	filters := Filters{
		Page:         1,
		PageSize:     20,
		Sort:         "id",
		SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
	}
	movies, _, err := testRepo.GetAll("Lord Of The Rings", []string{"Action"}, filters)
	if err != nil {
		t.Error("cannot get all movies")
	}

	if len(movies) != 1 {
		t.Errorf("expected 1 movie but got %d", len(movies))
	}
}

func TestPostgresDBRepoDeleteMovie(t *testing.T) {
	err := testRepo.Delete(movieID)
	if err != nil {
		t.Errorf("cannot delete movie: %s", err)
	}

	_, err = testRepo.Get(movieID)
	if err == nil {
		t.Error("got movie that was supposed to be deleted")
	}
}
