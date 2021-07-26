package transform

import (
	"strings"
)

// Method defines the method for the virtual types method
type Method struct {
	Name   string
	Action string
}

// Filter defines the data struct for the virtual types method
type Filter struct {
	// Scalar is the scalar name of a type, e.g. String, Int or DateTime
	Scalar  string
	Methods []Method
}

func (r *AST) filters() []Filter {
	var filters []Filter
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
		filters = append(filters, Filter{
			Scalar:  scalar,
			Methods: fields,
		})
	}
	return filters
}

// Filter returns a filter by scalar
func (r *AST) Filter(scalar string) *Filter {
	scalar = strings.Replace(scalar, "NullableFilter", "", 1)
	scalar = strings.Replace(scalar, "Filter", "", 1)
	for _, filter := range r.Filters {
		if filter.Scalar == scalar {
			return &filter
		}
	}
	return nil
}
