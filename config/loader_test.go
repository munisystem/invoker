package config

import (
	"reflect"
	"testing"
)

const fixtureDir = "fixture"

func TestLoadHcl(t *testing.T) {
	list, err := loadHcl(databaseRawString)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	expected := "*ast.ObjectList"
	actual := reflect.TypeOf(list).String()

	if expected != actual {
		t.Errorf("didn't match type: expected %s, actual %s", expected, actual)
	}
}

func TestLoadHcl_badFile(t *testing.T) {
	expected := "At 2:5: illegal char"
	_, actual := loadHcl(nopeRawString)
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

const databaseRawString = `
database "alice_db" {
  endpoint = "alice.example.com"
  port = 5432
  user = "admin"
  password = "admin"
  database_name = "apple"
}
`

const nopeRawString = `
nope:
  body: "It looks like JSON"
`
