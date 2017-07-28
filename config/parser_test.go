package config

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

func TestParseDatabaseConfig(t *testing.T) {
	obj, err := hcl.Parse(db)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList)

	expected := map[string]*Database{
		"alice_db": &Database{Endpoint: "alice.example.com", Port: 5432, User: "admin", Password: "admin", DatabaseName: "apple"},
	}

	actual, err := parseDatabaseConfig(list)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("didn't match struct: expected %v, actual %v", expected, actual)
	}
}

func TestParseDatabaseConfig_emptyName(t *testing.T) {
	obj, err := hcl.Parse(dbEmptyName)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList)

	expected := "2:1: database must be contained name"
	_, actual := parseDatabaseConfig(list)
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestParseDatabaseConfig_duplicateName(t *testing.T) {
	obj, err := hcl.Parse(dbDuplicateName)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList)

	expected := "3:1: alice_db is duplicate"
	_, actual := parseDatabaseConfig(list)
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestParseGroupConfig(t *testing.T) {
	obj, err := hcl.Parse(group)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList)

	expected := map[string]*Group{
		"dev": &Group{Policies: []string{"alice_db_readonly", "bob_db_readonly", "carol_db_writable"}},
	}

	actual, err := parseGroupConfig(list)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("didn't match struct: expected %v, actual %v", expected, actual)
	}
}

func TestParseGroupConfig_emptyName(t *testing.T) {
	obj, err := hcl.Parse(groupEmptyName)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList)

	expected := "2:1: group must be contained name"
	_, actual := parseGroupConfig(list)
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestParseGroupConfig_duplicateName(t *testing.T) {
	obj, err := hcl.Parse(groupDuplicateName)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList)

	expected := "3:1: dev is duplicate"
	_, actual := parseGroupConfig(list)
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

const db = `
database "alice_db" {
	endpoint = "alice.example.com"
	port = 5432
	user = "admin"
	password = "admin"
	database_name = "apple"
}
`

const dbEmptyName = `
database {}
`

const dbDuplicateName = `
database "alice_db" {}
database "alice_db" {}
`

const group = `
group "dev" {
  policies = [
    "alice_db_readonly",
    "bob_db_readonly",
    "carol_db_writable"
  ]
}
`

const groupEmptyName = `
group {}
`

const groupDuplicateName = `
group "dev" {}
group "dev" {}
`
