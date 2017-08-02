package postgresql

import (
	"database/sql"
	"fmt"
	"testing"

	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var pool *dockertest.Pool
var resource *dockertest.Resource
var db *sql.DB

func preparePostgreSQLContainer(t *testing.T) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		t.Fatalf("couldn't connect docker host: %s", err.Error())
	}

	resource, err = pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret"})
	if err != nil {
		t.Fatalf("couldn't start PostgreSQL container: %s", err.Error())
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatalf("couldn't prepare PostgreSQL container: %s", err.Error())
	}
}

func cleanupPostgreSQLContainer(t *testing.T) {
	if err := pool.Purge(resource); err != nil {
		t.Fatalf("couldn't cleanup PostgreSQL container: %s", err.Error())
	}
}

func TestRunQueries(t *testing.T) {
	preparePostgreSQLContainer(t)
	defer cleanupPostgreSQLContainer(t)

	pg := Initialize(fmt.Sprintf("postgres://postgres:secret@localhost:%s?sslmode=disable", resource.GetPort("5432/tcp")))
	if err := pg.RunQueries(TestQueries); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	expected := "Hello World!!"
	var actual string
	err := db.QueryRow("SELECT title FROM testing LIMIT 1;").Scan(&actual)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if expected != actual {
		t.Fatalf("didn't match be written string: expect %s, actual %s", expected, actual)
	}
}

var TestQueries = []string{
	"CREATE TABLE testing (title varchar(40));",
	"INSERT INTO testing VALUES ('Hello World!!')",
}
