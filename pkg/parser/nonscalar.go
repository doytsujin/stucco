package parser

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type namedDefintion interface {
	ast.Definition
	GetName() *ast.Name
}

func customDefinition(p *Parser, d ast.Definition) (gt graphql.Type, err error) {
	t, ok := d.(namedDefintion)
	if !ok {
		return nil, errors.New("not a type definition")
	}
	if gt, ok := p.gqlTypeMap[t.GetName().Value]; ok {
		return gt, nil
	}
	// Prevent recursion
	switch t := d.(type) {
	case *ast.ScalarDefinition:
		sc := graphql.ScalarConfig{
			Name: t.Name.Value,
		}
		if fn, ok := p.Scalars[t.Name.Value]; ok {
			sc.Serialize = fn.Serialize
			sc.ParseValue = fn.Parse
			sc.ParseLiteral = func(v ast.Value) interface{} {
				return fn.Parse(v.GetValue())
			}
		}
		setDescription(&sc.Description, t)
		st := graphql.NewScalar(sc)
		p.gqlTypeMap[st.Name()] = st
		gt = st
	case *ast.EnumDefinition:
		gt, err = enumDefinition(p, t)
	case *ast.InputObjectDefinition:
		gt, err = inputObjectDefinition(p, t)
	case *ast.InterfaceDefinition:
		gt, err = interfaceDefinition(p, t)
	case *ast.ObjectDefinition:
		gt, err = objectDefintion(p, t)
	case *ast.UnionDefinition:
		gt, err = unionDefinition(p, t)
	default:
		err = errors.New("unsupported type defintion")
	}
	return
}
