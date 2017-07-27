package config

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseDatabaseConfig(t *testing.T) {
	path := filepath.Join(fixtureDir, "database.hcl")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list, err := loadHcl(string(b))
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	expected, err := parseDatabaseConfig(list)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	actual := map[string]*Database{
		"alice_db": &Database{Endpoint: "alice.example.com", Port: 5432, User: "admin", Password: "admin", DatabaseName: "apple"},
		"bob_db":   &Database{Endpoint: "bob.example.com", Port: 5432, User: "admin", Password: "admin", DatabaseName: "banana"},
		"carol_db": &Database{Endpoint: "carol.example.com", Port: 5432, User: "admin", Password: "admin", DatabaseName: "cherry"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("didn't match struct: expected %v, actual %v", expected, actual)
	}
}
