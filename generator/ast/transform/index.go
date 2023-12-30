package transform

import (
	"github.com/steebchen/prisma-client-go/generator/ast/dmmf"
	"github.com/steebchen/prisma-client-go/generator/types"
)

type Index struct {
	Name         types.String   `json:"name"`
	InternalName string         `json:"internalName"`
	Fields       []types.String `json:"fields"`
}

func indexes(m dmmf.Model) []Index {
	var idx []Index
	for _, i := range m.UniqueIndexes {
		internalName := i.InternalName
		if internalName == "" {
			internalName = concatFieldsToName(i.Fields)
		}
		idx = append(idx, Index{
			Name:         getName(i.InternalName, i.Fields),
			InternalName: internalName,
			Fields:       i.Fields,
		})
	}

	if len(m.PrimaryKey.Fields) > 0 {
		idx = append(idx, Index{
			Name:         getName(concatFieldsToName(m.PrimaryKey.Fields), m.PrimaryKey.Fields),
			InternalName: concatFieldsToName(m.PrimaryKey.Fields),
			Fields:       m.PrimaryKey.Fields,
		})
	}

	return idx
}

func concatFieldsToName(fields []types.String) string {
	var name string
	for i, f := range fields {
		if i > 0 {
			name += "_"
		}
		name += f.String()
	}
	return name
}

func getName(field string, fields []types.String) types.String {
	if field != "" {
		return types.String(field)
	}
	var name string
	for _, f := range fields {
		name += f.GoCase()
	}
	return types.String(name)
}
