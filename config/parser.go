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
		if len(item.Keys) != 2 {
			return nil, fmt.Errorf("%s: database must be contained name", item.Pos())
		}

		var d Database
		if err := hcl.DecodeObject(&d, item.Val); err != nil {
			return nil, err
		}

		name := item.Keys[1].Token.Value().(string)
		if _, exists := databases[name]; exists {
			return nil, fmt.Errorf("%s: %s is duplicate", item.Pos(), name)
		}

		databases[name] = &d
	}

	return databases, nil
}
