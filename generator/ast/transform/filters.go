package transform

// Method defines the method for the virtual types method
type Method struct {
	Name   string
	Action string
}

// Type defines the data struct for the virtual types method
type Type struct {
	Name    string
	Methods []Method
}

func (r *AST) filters() []Type {
	var filters []Type
	for _, scalar := range r.Scalars {
		p := r.pick(scalar + "Filter")
		if p == nil {
			p = r.pick(scalar + "NullableFilter")
			if p == nil {
				continue
			}
		}
		var fields []Method
		for _, field := range p.Fields {
			// specifically ignore equals, as it gets special handling
			if field.Name == "equals" {
				continue
			}
			fields = append(fields, Method{
				Name:   field.Name.GoCase(),
				Action: field.Name.String(),
			})
		}
		filters = append(filters, Type{
			Name:    scalar,
			Methods: fields,
		})
	}
	return filters
}
