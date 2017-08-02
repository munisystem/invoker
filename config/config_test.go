package config

import "testing"

func TestCheckDependencies(t *testing.T) {
	config := &Config{
		Databases: map[string]*Database{
			"alice_db": &Database{},
		},
		Policies: map[string]*Policy{
			"alice_db_readonly": &Policy{
				Database: "alice_db",
			},
		},
		Groups: map[string]*Group{
			"dev": &Group{
				Policies: []string{
					"alice_db_readonly",
				},
			},
		},
		Users: map[string]*User{
			"alice": &User{
				Group: "dev",
			},
		},
	}

	if err := config.CheckDependencies(); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}
}

func TestCheckDependencies_missingGroup(t *testing.T) {
	config := &Config{
		Groups: map[string]*Group{
			"dev":  &Group{},
			"core": &Group{},
		},
		Users: map[string]*User{
			"alice": &User{
				Group: "dev",
			},
			"bob": &User{
				Group: "core",
			},
			"carol": &User{
				Group: "owner",
			},
		},
	}

	expected := "didn't declare group of 'owner'"
	if actual := config.CheckDependencies(); actual != nil && expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestCheckDependencies_missingPolicy(t *testing.T) {
	config := &Config{
		Policies: map[string]*Policy{
			"alice_db_readonly": &Policy{},
			"alice_db_writable": &Policy{},
			"bob_db_owner":      &Policy{},
		},
		Groups: map[string]*Group{
			"dev": &Group{
				Policies: []string{
					"alice_db_readonly",
					"bob_db_writable",
				},
			},
			"core": &Group{
				Policies: []string{
					"alice_db_writable",
				},
			},
		},
	}

	expected := "didn't declare policy of 'bob_db_writable'"
	if actual := config.CheckDependencies(); actual != nil && expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}

func TestCheckDependencies_missingDatabase(t *testing.T) {
	config := &Config{
		Databases: map[string]*Database{
			"alice_db": &Database{},
			"bob_db":   &Database{},
		},
		Policies: map[string]*Policy{
			"alice_db_readonly": &Policy{
				Database: "alice_db",
			},
			"bob_db_readonly": &Policy{
				Database: "bob_db",
			},

			"carol_db_writable": &Policy{
				Database: "carol_db",
			},
		},
	}

	expected := "didn't declare database of 'carol_db'"
	if actual := config.CheckDependencies(); actual != nil && expected != actual.Error() {
		t.Fatalf("didn't match err: expected %s, actual %s", expected, actual.Error())
	}
}
