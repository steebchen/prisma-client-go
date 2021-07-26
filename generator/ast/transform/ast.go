package transform

import (
	"fmt"
	"github.com/prisma/prisma-client-go/generator/ast/dmmf"
)

type AST struct {
	dmmf *dmmf.Document

	// Scalars describe a list of scalar types, such as Int, String, DateTime, etc.
	Scalars []string `json:"scalars"`
	// Filters describe a list of scalar types and the respective read operations
	Filters []Filter `json:"filters"`
}

func New(document *dmmf.Document) *AST {
	ast := &AST{
		dmmf: document,
	}

	ast.Scalars = ast.scalars()
	ast.Filters = ast.filters()

	return ast
}

func (r *AST) pick(name string) *dmmf.CoreType {
	for _, i := range r.dmmf.Schema.InputObjectTypes.Prisma {
		if string(i.Name) == name {
			return &i
		}
	}
	fmt.Printf("no type %s found\n", name)
	return nil
}
