package transform

import (
	"github.com/steebchen/prisma-client-go/generator/ast/dmmf"
)

type AST struct {
	dmmf *dmmf.Document

	// Scalars describe a list of scalar types, such as Int, String, DateTime, etc.
	Scalars []string `json:"scalars"`

	// Enums contains all enums
	Enums []Enum `json:"enums"`

	// Models contains top-level information including fields and their respective filters
	Models []Model `json:"models"`

	// ReadFilters describe a list of scalar types and the respective read operations
	ReadFilters []Filter `json:"readFilters"`

	// WriteFilters describe a list of scalar types and the respective read operations
	WriteFilters []Filter `json:"writeFilters"`

	// OrderBys describe a list of what FindMany operations can order by
	OrderBys []OrderBy `json:"orderBys"`
}

func New(document *dmmf.Document) *AST {
	ast := &AST{
		dmmf: document,
	}

	// first, fetch types
	ast.Scalars = ast.scalars()
	ast.Enums = ast.enums()

	// fetch data
	ast.Models = ast.models()

	// fetch data which is needed for the query api, which require ast types
	ast.ReadFilters = ast.readFilters()
	ast.WriteFilters = ast.writeFilters()

	// add old, deprecated filters which are just added for compatibility reasons
	// these can be removed at some point
	for _, filter := range ast.deprecatedReadFilters() {
		for i, f := range ast.ReadFilters {
			if f.Name == filter.Name {
				ast.ReadFilters[i].Methods = append(ast.ReadFilters[i].Methods, filter.Methods...)
			}
		}
	}

	return ast
}

func (r *AST) pick(names ...string) *dmmf.CoreType {
	for _, name := range names {
		for _, i := range r.dmmf.Schema.InputObjectTypes.Prisma {
			if string(i.Name) == name {
				return &i
			}
		}
	}
	return nil
}
