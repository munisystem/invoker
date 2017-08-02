package postgresql

import (
	"database/sql"
	"fmt"
	"testing"

	dockertest "gopkg.in/ory-am/dockertest.v3"
)

func preparePostgreSQLContainer(t *testing.T) (func(), string, *sql.DB) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("couldn't connect docker host: %s", err.Error())
	}

	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret"})
	if err != nil {
		t.Fatalf("couldn't start PostgreSQL container: %s", err.Error())
	}

	addr := fmt.Sprintf("postgres://postgres:secret@localhost:%s?sslmode=disable", resource.GetPort("5432/tcp"))

	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("couldn't cleanup PostgreSQL container: %s", err.Error())
		}
	}

	var db *sql.DB
	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", addr)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatalf("couldn't prepare PostgreSQL container: %s", err.Error())
	}

	return cleanup, addr, db
}

func TestRunQueries(t *testing.T) {
	cleanup, addr, db := preparePostgreSQLContainer(t)
	defer cleanup()

	pg := Initialize(addr)
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
