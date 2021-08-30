package transform

func (r *AST) writeFilters() []Filter {
	var filters []Filter
	for _, scalar := range r.Scalars {
		p := r.pick(scalar + "FieldUpdateOperationsInput")
		if p == nil {
			p = r.pick("Nullable" + scalar + "FieldUpdateOperationsInput")
			if p == nil {
				continue
			}
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
	return filters
}

// WriteFilter returns a filter for a read operation by scalar
func (r *AST) WriteFilter(scalar string) *Filter {
	for _, filter := range r.WriteFilters {
		if filter.Name == scalar {
			return &filter
		}
	}
	return nil
}
