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

	list := obj.Node.(*ast.ObjectList).Filter("database")

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

	list := obj.Node.(*ast.ObjectList).Filter("database")

	expected := "2:10: database must be contained name"
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

	list := obj.Node.(*ast.ObjectList).Filter("database")

	expected := "3:21: alice_db is duplicate"
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

	list := obj.Node.(*ast.ObjectList).Filter("group")

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

	list := obj.Node.(*ast.ObjectList).Filter("group")

	expected := "2:7: group must be contained name"
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

	list := obj.Node.(*ast.ObjectList).Filter("group")

	expected := "3:13: dev is duplicate"
	_, actual := parseGroupConfig(list)
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestPolicyGroupConfig(t *testing.T) {
	obj, err := hcl.Parse(policy)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList).Filter("policy")

	expected := map[string]*Policy{
		"alice_db_writable": &Policy{
			Database: "alice_db",
			Queries: []string{
				"CREATE ROLE {{ .Name }} WITH LOGIN ENCRYPTED PASSWORD {{ .Password }};",
				"GRANT ALL ON ALL TABLES IN SCHEMA public TO {{ .Name }};",
				"GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO {{ .Name }};",
			},
		},
	}

	actual, err := parsePolicyConfig(list)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("didn't match struct: expected %v, actual %v", expected, actual)
	}
}

func TestParsePolicyConfig_emptyName(t *testing.T) {
	obj, err := hcl.Parse(policyEmptyName)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList).Filter("policy")

	expected := "2:8: policy must be contained name"
	_, actual := parsePolicyConfig(list)
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestParsePolicyConfig_duplicateName(t *testing.T) {
	obj, err := hcl.Parse(policyDuplicateName)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList).Filter("policy")

	expected := "3:28: alice_db_writable is duplicate"
	_, actual := parsePolicyConfig(list)
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestParseUserConfig(t *testing.T) {
	obj, err := hcl.Parse(user)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList).Filter("user")

	expected := map[string]*User{
		"cat": &User{Group: "core"},
	}

	actual, err := parseUserConfig(list)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("didn't match struct: expected %v, actual %v", expected, actual)
	}
}

func TestParseUserConfig_emptyName(t *testing.T) {
	obj, err := hcl.Parse(userEmptyName)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList).Filter("user")

	expected := "2:6: user must be contained name"
	_, actual := parseUserConfig(list)
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestParseUserConfig_duplicateName(t *testing.T) {
	obj, err := hcl.Parse(userDuplicateName)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list := obj.Node.(*ast.ObjectList).Filter("user")

	expected := "3:12: cat is duplicate"
	_, actual := parseUserConfig(list)
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

const policy = `
policy "alice_db_writable" {
  database = "alice_db"

  queries = [
    "CREATE ROLE {{ .Name }} WITH LOGIN ENCRYPTED PASSWORD {{ .Password }};",
    "GRANT ALL ON ALL TABLES IN SCHEMA public TO {{ .Name }};",
    "GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO {{ .Name }};"
  ]
}
`

const policyEmptyName = `
policy {}
`

const policyDuplicateName = `
policy "alice_db_writable" {}
policy "alice_db_writable" {}
`

const user = `
user "cat" {
	group = "core"
}
`

const userEmptyName = `
user {}
`

const userDuplicateName = `
user "cat" {}
user "cat" {}
`
