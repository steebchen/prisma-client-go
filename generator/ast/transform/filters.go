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

func (r *AST) readFilters() []Filter {
	var filters []Filter
	for _, scalar := range r.Scalars {
		p := r.pick(scalar + "ReadFilter")
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

// ReadFilter returns a filter for a read operation by scalar
func (r *AST) ReadFilter(scalar string) *Filter {
	scalar = strings.Replace(scalar, "NullableFilter", "", 1)
	scalar = strings.Replace(scalar, "ReadFilter", "", 1)
	for _, filter := range r.ReadFilters {
		if filter.Scalar == scalar {
			return &filter
		}
	}
	return nil
}
