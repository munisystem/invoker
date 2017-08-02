package config

import "fmt"

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

func (c *Config) CheckDependencies() error {
	// Chack User dependencies
	for _, user := range c.Users {
		exists := false
		for group := range c.Groups {
			if user.Group == group {
				exists = true
				break
			}
		}
		if !exists {
			return fmt.Errorf("didn't declare group of '%s'", user.Group)
		}
	}

	// Chack Group dependencies
	for _, group := range c.Groups {

		// Group has some policies
		for _, gp := range group.Policies {
			exists := false
			for policy := range c.Policies {
				if gp == policy {
					exists = true
					break
				}
			}

			// Chack Group having policy is declare
			if !exists {
				return fmt.Errorf("didn't declare policy of '%s'", gp)
			}
		}
	}

	// Check Policy dependencies
	for _, policy := range c.Policies {
		exists := false
		for database := range c.Databases {
			if policy.Database == database {
				exists = true
				break
			}
		}
		if !exists {
			return fmt.Errorf("didn't declare database of '%s'", policy.Database)
		}
	}

	return nil
}
