package config

import (
	"path/filepath"
	"reflect"
	"testing"
)

const fixtureDir = "fixture"

func TestLoadFile(t *testing.T) {
	path := filepath.Join(fixtureDir, "mixed_keys.hcl")

	expected := &Config{
		Databases: map[string]*Database{
			"alice_db": &Database{
				Endpoint:     "alice.example.com",
				Port:         5432,
				User:         "admin",
				Password:     "admin",
				DatabaseName: "apple",
			},
		},
		Groups: map[string]*Group{
			"dev": &Group{
				Policies: []string{"alice_db_readonly"},
			},
		},
		Policies: map[string]*Policy{
			"alice_db_readonly": &Policy{
				Database: "alice_db",
				Queries: []string{
					"CREATE ROLE {{ .Name }} WITH LOGIN ENCRYPTED PASSWORD {{ .Password }};",
					"GRANT SELECT ON ALL TABLES IN SCHEMA public TO {{ .Name }};",
					"GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO {{ .Name }};",
				},
			},
		},
		Users: map[string]*User{
			"alice": &User{
				Group: "dev",
			},
		},
	}

	actual, err := LoadFile(path)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("didn't match struct: expected %v, actual %v", expected, actual)
	}
}

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
