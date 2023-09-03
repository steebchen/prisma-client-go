package transform

import (
	"github.com/steebchen/prisma-client-go/generator/ast/dmmf"
	"github.com/steebchen/prisma-client-go/generator/types"
)

type Model struct {
	Name    types.String `json:"name"`
	Fields  []Field      `json:"fields"`
	Indexes []Index      `json:"indexes"`

	// TODO remove this and apply all required data directly to model
	OldModel dmmf.Model `json:"-"`
}

func (m Model) CompoundKeys() []Index {
	var items []Index
	items = append(items, m.Indexes...)

	if m.OldModel.PrimaryKey.Name != "" {
		items = append(items, Index{
			Name:         m.OldModel.PrimaryKey.Name,
			InternalName: m.OldModel.PrimaryKey.Name.String(),
			Fields:       m.OldModel.PrimaryKey.Fields,
		})
	}

	return items
}

type Field struct {
	// TODO re-declare all fields here instead of embedding dmmf.Field

	// Prisma indicates whether this is a pseudo field used for Prisma-specific actions, e.g. 'Relevance_'
	Prisma bool `json:"prisma"`

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
