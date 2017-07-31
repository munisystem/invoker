package config

import (
	"errors"
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
