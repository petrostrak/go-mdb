package data

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
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
	code := m.Run()

	os.Exit(code)
}
