package transform

import (
	"github.com/prisma/prisma-client-go/generator/ast/dmmf"
)

// applyDeprecatedRelationReadFilters applies old relation read methods, the only one being `Where` -> `Is`
func (r *AST) applyDeprecatedRelationReadFilters() {
	for m, model := range r.Models {
		for f, field := range model.Fields {
			for _, relationField := range field.RelationReadMethods {
				if relationField.Name != "Is" {
					continue
				}
				// for the relation field method which is the operation `Is`, append the same action but of name `Where`
				r.Models[m].Fields[f].RelationReadMethods = append(r.Models[m].Fields[f].RelationReadMethods, Method{
					Name:       "Where",
					Action:     "is",
					Deprecated: "Is",
				})
			}
		}
	}
}

func (r *AST) applyRelationReadFilters() {
	var foundField dmmf.OuterInputType

	for _, model := range r.Models {
		combinations := []string{
			model.Name.String() + "RelationFilter",
			model.Name.String() + "ListRelationFilter",
		}

		for _, name := range combinations {
			// get <Model>WhereInput
			// iterate through <<Model>WhereInput>.fields as XField
		outer:
			for _, innerModel := range r.Models {
				p := r.pick(innerModel.Name.String() + "WhereInput")
				if p == nil {
					continue
				}
				for _, field := range p.Fields {
					// iterate through <XField>.InputTypes
					for _, input := range field.InputTypes {
						// check for inputType.type == <Model>(List)RelationFilter
						if input.Type.String() == name {
							foundField = field
							break outer
						}
					}
				}
			}

			// need to iterate through models again as fields reference each other in different models
			for m, innerModel := range r.Models {
				item := r.pick(name)
				if item == nil {
					continue
				}

				for f, field := range innerModel.Fields {
					if field.Name == foundField.Name {
						for _, relationField := range item.Fields {
							// TODO prevent duplicate items. this should be refactored
							var has bool
							for _, method := range r.Models[m].Fields[f].RelationReadMethods {
								if method.Name == relationField.Name.GoCase() {
									has = true
								}
							}
							if !has {
								r.Models[m].Fields[f].RelationReadMethods = append(r.Models[m].Fields[f].RelationReadMethods, Method{
									Name:   relationField.Name.GoCase(),
									Action: relationField.Name.String(),
									// TODO – this is currently used from the field relation, but would be more useful to provide here
									//IsList: itemField.IsList,
								})
							}
						}
					}
				}
			}
		}
	}
}
