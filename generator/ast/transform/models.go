package transform

import (
	"github.com/prisma/prisma-client-go/generator/ast/dmmf"
	"github.com/prisma/prisma-client-go/generator/types"
)

type Model struct {
	Name    types.String `json:"name"`
	Fields  []Field      `json:"fields"`
	Indexes []Index      `json:"indexes"`

	// TODO remove this and apply all required data directly to model
	OldModel dmmf.Model
}

type Field struct {
	// TODO re-declare all fields here instead of embedding dmmf.Field

	dmmf.Field
}

func (r *AST) models() []Model {
	var models []Model
	for _, model := range r.dmmf.Datamodel.Models {
		var fields []Field
		for _, field := range model.Fields {
			fields = append(fields, Field{
				Field: field,
			})
		}
		m := Model{
			Name:     model.Name,
			Fields:   fields,
			OldModel: model,
		}
		m.Indexes = indexes(model)
		models = append(models, m)
	}
	return models
}
