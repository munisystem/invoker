package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

func parseDatabaseConfig(list *ast.ObjectList) (map[string]*Database, error) {
	if len(list.Items) == 0 {
		return nil, nil
	}

	databases := make(map[string]*Database)
	for _, item := range list.Items {
		if len(item.Keys) != 1 {
			return nil, fmt.Errorf("%s: database must be contained name", item.Val.Pos())
		}

		var d Database
		if err := hcl.DecodeObject(&d, item.Val); err != nil {
			return nil, err
		}

		name := item.Keys[0].Token.Value().(string)
		if _, exists := databases[name]; exists {
			return nil, fmt.Errorf("%s: %s is duplicate", item.Val.Pos(), name)
		}

		databases[name] = &d
	}

	return databases, nil
}

func parseGroupConfig(list *ast.ObjectList) (map[string]*Group, error) {
	if len(list.Items) == 0 {
		return nil, nil
	}

	groups := make(map[string]*Group)
	for _, item := range list.Items {
		if len(item.Keys) != 1 {
			return nil, fmt.Errorf("%s: group must be contained name", item.Val.Pos())
		}

		var g Group
		if err := hcl.DecodeObject(&g, item.Val); err != nil {
			return nil, err
		}

		name := item.Keys[0].Token.Value().(string)
		if _, exists := groups[name]; exists {
			return nil, fmt.Errorf("%s: %s is duplicate", item.Val.Pos(), name)
		}

		groups[name] = &g
	}

	return groups, nil
}

func parsePolicyConfig(list *ast.ObjectList) (map[string]*Policy, error) {
	if len(list.Items) == 0 {
		return nil, nil
	}

	policies := make(map[string]*Policy)
	for _, item := range list.Items {
		if len(item.Keys) != 1 {
			return nil, fmt.Errorf("%s: policy must be contained name", item.Val.Pos())
		}

		var p Policy
		if err := hcl.DecodeObject(&p, item.Val); err != nil {
			return nil, err
		}

		name := item.Keys[0].Token.Value().(string)
		if _, exists := policies[name]; exists {
			return nil, fmt.Errorf("%s: %s is duplicate", item.Val.Pos(), name)
		}

		policies[name] = &p
	}

	return policies, nil
}

func parseUserConfig(list *ast.ObjectList) (map[string]*User, error) {
	if len(list.Items) == 0 {
		return nil, nil
	}

	users := make(map[string]*User)
	for _, item := range list.Items {
		if len(item.Keys) != 1 {
			return nil, fmt.Errorf("%s: user must be contained name", item.Val.Pos())
		}

		var u User
		if err := hcl.DecodeObject(&u, item.Val); err != nil {
			return nil, err
		}

		name := item.Keys[0].Token.Value().(string)
		if _, exists := users[name]; exists {
			return nil, fmt.Errorf("%s: %s is duplicate", item.Val.Pos(), name)
		}

		users[name] = &u
	}

	return users, nil
}
