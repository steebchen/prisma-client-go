package transform

import (
	"github.com/steebchen/prisma-client-go/generator/types"
)

type Enum struct {
	Name   types.String   `json:"name"`
	Values []types.String `json:"values"`
}

func (r *AST) enums() []Enum {
	var enums []Enum
	for _, enum := range r.dmmf.Schema.EnumTypes.Model {
		enums = append(enums, Enum{
			Name:   enum.Name,
			Values: enum.Values,
		})
	}
	return enums
}
