package config

type Config struct {
	Databases map[string]*Database
	Groups    map[string]*Group
	Policies  map[string]*Policy
	Users     map[string]*User
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
	Database string   `hcl:"database"`
	Queries  []string `hcl:"queries"`
}

type User struct {
	Group string `hcl:"group"`
}

func NewConfig() *Config {
	return &Config{
		Databases: make(map[string]*Database),
		Groups:    make(map[string]*Group),
		Policies:  make(map[string]*Policy),
		Users:     make(map[string]*User),
	}
}
