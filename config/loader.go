package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

func LoadFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	list, err := loadHcl(string(b))
	if err != nil {
		return nil, err
	}

	config := new(Config)

	if f := list.Filter("database"); len(f.Items) > 0 {
		dbs, err := parseDatabaseConfig(f)
		if err != nil {
			return nil, err
		}
		config.Databases = dbs
	}

	if f := list.Filter("group"); len(f.Items) > 0 {
		groups, err := parseGroupConfig(f)
		if err != nil {
			return nil, err
		}
		config.Groups = groups
	}

	if f := list.Filter("policy"); len(f.Items) > 0 {
		policies, err := parsePolicyConfig(f)

		if err != nil {
			return nil, err
		}
		config.Policies = policies
	}

	if f := list.Filter("user"); len(f.Items) > 0 {
		users, err := parseUserConfig(f)

		if err != nil {
			return nil, err
		}
		config.Users = users
	}

	return config, nil
}

func loadHcl(hclText string) (*ast.ObjectList, error) {
	obj, err := hcl.Parse(hclText)
	if err != nil {
		return nil, err
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, errors.New("failed to parse: file does not contain root object")
	}

	return list, nil
}

func mergeConfig(dst, src *Config) error {
	if len(dst.Databases) > 0 {
		for k, v := range dst.Databases {
			if _, ok := src.Databases[k]; ok {
				return fmt.Errorf("database '%s' is duplicate", k)
			}
			src.Databases[k] = v
		}
	}

	if len(dst.Groups) > 0 {
		for k, v := range dst.Groups {
			if _, ok := src.Groups[k]; ok {
				return fmt.Errorf("group '%s' is duplicate", k)
			}
			src.Groups[k] = v
		}
	}

	if len(dst.Policies) > 0 {
		for k, v := range dst.Policies {
			if _, ok := src.Policies[k]; ok {
				return fmt.Errorf("policy '%s' is duplicate", k)
			}
			src.Policies[k] = v
		}
	}

	if len(dst.Users) > 0 {
		for k, v := range dst.Users {
			if _, ok := src.Users[k]; ok {
				return fmt.Errorf("users '%s' is duplicate", k)
			}
			src.Users[k] = v
		}
	}
	return nil
}
