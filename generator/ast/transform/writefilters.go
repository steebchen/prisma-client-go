package transform

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
			fields = append(fields, Method{
				Name:   field.Name.GoCase(),
				Action: field.Name.String(),
			})
		}
		filters = append(filters, Filter{
			Name:    scalar,
			Methods: fields,
		})
	}
	// pick out extra methods, e.g. for list write operations
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
						if inputType.Location == "scalar" {
							scalarName = inputType.Type.String() + "List" // create an on-the-fly <scalar>List filter type
						}
					}
					continue
				}
				fields = append(fields, Method{
					Name:   field.Name.GoCase(),
					Action: field.Name.String(),
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
