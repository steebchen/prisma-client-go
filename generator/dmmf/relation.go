package dmmf

import (
	"github.com/prisma/prisma-client-go/generator/types"
)

// ReverseRelationName returns a relation name of a given field name.
// For example, when passing the types from=User and to=Post, the function may return "author". For the other way round,
// it may return "posts".
func (d *Document) ReverseRelationName(from types.StringLike, to types.StringLike) types.String {
	for _, output := range d.Schema.OutputTypes {
		if output.Name.String() != from.String() {
			continue
		}

		for _, field := range output.Fields {
			if field.OutputType.Type.String() != to.String() {
				continue
			}

			return field.Name
		}
	}

	return ""
}
