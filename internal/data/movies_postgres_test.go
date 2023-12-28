package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host   = "localhost"
	user   = "postgres"
	pw     = "password"
	dbName = "movie_test"
	port   = "5435"
	dns    = "host=%s port=%s user=%s password=%s dbname=%s sslmod=disabled timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}

	pool = p

	opts := dockertest.RunOptions{
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

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		pool.Purge(resource)
		log.Fatal(err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pq", fmt.Sprintf(dns, host, port, user, pw, dbName))
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("could not connect to db: %s", err)
	}

	code := m.Run()

	os.Exit(code)
}
