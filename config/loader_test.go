package config

import (
	"path/filepath"
	"reflect"
	"testing"
)

const fixtureDir = "fixture"

func TestLoadDir(t *testing.T) {
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

	actual, err := LoadDir(fixtureDir)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("didn't match struct: expected %v, actual %v", expected, actual)
	}
}

func TestLoadDir_duplicateName(t *testing.T) {
	expected := "database 'alice_db' is duplicate"
	if _, actual := LoadDir(filepath.Join(fixtureDir, "duplicate")); expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestLoadFile(t *testing.T) {
	path := filepath.Join(fixtureDir, "database.hcl")

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

func TestMergeConfig(t *testing.T) {
	aliceDB := &Database{
		Endpoint:     "alice.example.com",
		Port:         5432,
		User:         "admin",
		Password:     "admin",
		DatabaseName: "apple",
	}

	bobDB := &Database{
		Endpoint:     "bob.example.com",
		Port:         5432,
		User:         "admin",
		Password:     "admin",
		DatabaseName: "banana",
	}

	dst := NewConfig()
	dst.Databases["alice_db"] = aliceDB

	actual := NewConfig()
	actual.Databases["bob_db"] = bobDB

	expected := NewConfig()
	expected.Databases["alice_db"] = aliceDB
	expected.Databases["bob_db"] = bobDB

	if err := mergeConfig(dst, actual); err != nil {
		t.Errorf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("didn't match struct: expected %v, actual %v", expected, actual)
	}
}

func TestMergeConfig_DuplicateName(t *testing.T) {
	aliceDB := &Database{
		Endpoint:     "alice.example.com",
		Port:         5432,
		User:         "admin",
		Password:     "admin",
		DatabaseName: "apple",
	}

	dst := NewConfig()
	dst.Databases["alice_db"] = aliceDB

	src := NewConfig()
	src.Databases["alice_db"] = aliceDB

	expected := "database 'alice_db' is duplicate"
	if actual := mergeConfig(dst, src); expected != actual.Error() {
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
