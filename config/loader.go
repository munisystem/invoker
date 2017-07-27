package config

import (
	"errors"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

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
