package config

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

const fixtureDir = "fixture"

func TestLoadHcl(t *testing.T) {
	path := filepath.Join(fixtureDir, "database.hcl")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	list, err := loadHcl(string(b))
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
	path := filepath.Join(fixtureDir, "nope.hcl")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	expected := "At 1:5: illegal char"
	_, actual := loadHcl(string(b))
	if expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}
