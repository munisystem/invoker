package config

type Config struct {
	Databases map[string]*Database
	Groups    map[string]*Group
	Policies  []*Policy
	Users     []*User
}

type Database struct {
	Endpoint     string `hcl:"endpoint"`
	Port         int    `hcl:"port"`
	User         string `hcl:"user"`
	Password     string `hcl:"password"`
	DatabaseName string `hcl:"database_name"`
}

type Group struct {
	Policies []string `hcl:"policies"`
}

type Policy struct {
	Name string `hcl:"-"`

	Database string   `hcl:"database"`
	Queries  []string `hcl:"queries"`
}

type User struct {
	Name string `hcl:"-"`

	Group string `hcl:"group"`
}
