package transform

import (
	"github.com/steebchen/prisma-client-go/generator/types"
)

func (r *AST) writeFilters() []Filter {
	var filters []Filter
	for _, scalar := range r.Scalars {
		p := r.pick(
			scalar+"FieldUpdateOperationsInput",
			"Nullable"+scalar+"FieldUpdateOperationsInput",
		)
		if p == nil {
			continue
		}
		var fields []Method
		for _, field := range p.Fields {
			// specifically ignore equals, as it gets special handling
			if field.Name == "set" {
				continue
			}
			var typeName types.Type
			var isList bool
			for _, inputType := range field.InputTypes {
				if inputType.Location == "scalar" && inputType.Type != "Null" {
					typeName = inputType.Type
					if inputType.IsList {
						isList = true
					}
				}
			}
			fields = append(fields, Method{
				Name:   field.Name.GoCase(),
				Action: field.Name.String(),
				Type:   typeName,
				IsList: isList,
			})
		}
		filters = append(filters, Filter{
			Name:    scalar,
			Methods: fields,
		})
	}
	// pick out extra methods, e.g. for list write operations
	// this could all be removed if the DMMF would implement operations by scalar, e.g. similar to
	// StringFilter there could be a write operation of StringListWriteOperations instead of individual
	// dmmf for each model+field combination
	for _, model := range r.Models {
		for _, field := range model.Fields {
			p := r.pick(
				model.Name.String() + "Update" + field.Name.String() + "Input",
			)
			if p == nil {
				continue
			}
			var scalarName string
			var fields []Method
			for _, field := range p.Fields {
				// specifically ignore equals, as it gets special handling
				if field.Name == "set" {
					for _, inputType := range field.InputTypes {
						if inputType.Location == "scalar" && inputType.Type != "Null" {
							scalarName = inputType.Type.String() + "List" // create an on-the-fly <scalar>List filter type
						}
					}
					continue
				}
				var typeName types.Type
				var isList bool
				for _, inputType := range field.InputTypes {
					if inputType.Location == "scalar" && inputType.Type != "Null" {
						typeName = inputType.Type
						if inputType.IsList {
							isList = true
						}
					}
				}
				fields = append(fields, Method{
					Name:   field.Name.GoCase(),
					Action: field.Name.String(),
					Type:   typeName,
					IsList: isList,
				})
			}
			filters = append(filters, Filter{
				Name:    scalarName,
				Methods: fields,
			})
		}
	}
	return filters
}

// WriteFilter returns a filter for a read operation by scalar
func (r *AST) WriteFilter(scalar string, isList bool) *Filter {
	if isList {
		scalar += "List"
	}
	for _, filter := range r.WriteFilters {
		if filter.Name == scalar {
			return &filter
		}
	}
	return nil
}
