package transform

import (
	"github.com/prisma/prisma-client-go/generator/ast/dmmf"
	"strings"
)

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
			if method := convertField(field); method != nil {
				fields = append(fields, *method)
			}
		}
		filters = append(filters, Filter{
			Name:    scalar,
			Methods: fields,
		})
	}
	for _, enum := range r.Enums {
		p := r.pick("Enum" + enum.Name + "Filter")
		if p == nil {
			p = r.pick("Enum" + enum.Name + "NullableFilter")
			if p == nil {
				continue
			}
		}

		var fields []Method
		for _, field := range p.Fields {
			if method := convertField(field); method != nil {
				fields = append(fields, *method)
			}
		}

		filters = append(filters, Filter{
			Name:    enum.Name,
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
		if filter.Name == scalar {
			return &filter
		}
	}
	return nil
}

func convertField(field dmmf.OuterInputType) *Method {
	// specifically ignore equals, as it gets special handling
	if field.Name == "equals" {
		return nil
	}
	isList := false
	// check if any of the input types accept a list
	for _, x := range field.InputTypes {
		// if yes, consider it a list item, regardless of all other items
		if x.IsList {
			isList = true
		}
	}
	return &Method{
		Name:   field.Name.GoCase(),
		Action: field.Name.String(),
		IsList: isList,
	}
}

// deprecatedReadFilters contains a hard-coded list of old filters to not breaking existing users
// these can be removed at some point in the future
func (r *AST) deprecatedReadFilters() []Filter {
	numberFilters := []Method{
		{
			Name:       "LT",
			Action:     "lt",
			Deprecated: "Lt",
		},
		{
			Name:       "LTE",
			Action:     "lte",
			Deprecated: "Lte",
		},
		{
			Name:       "GT",
			Action:     "gt",
			Deprecated: "Gt",
		},
		{
			Name:       "GTE",
			Action:     "gte",
			Deprecated: "Gte",
		},
	}
	return []Filter{
		{
			Name:    "Int",
			Methods: numberFilters,
		},
		{
			Name:    "Float",
			Methods: numberFilters,
		},
		{
			Name: "String",
			Methods: []Method{
				{
					Name:       "HasPrefix",
					Action:     "starts_with",
					Deprecated: "StartsWith",
				},
				{
					Name:       "HasSuffix",
					Action:     "ends_with",
					Deprecated: "EndsWith",
				},
			},
		},
		{
			Name: "DateTime",
			Methods: []Method{
				{
					Name:       "Before",
					Action:     "lt",
					Deprecated: "Lt",
				},
				{
					Name:       "After",
					Action:     "gt",
					Deprecated: "Gt",
				},
				{
					Name:       "BeforeEquals",
					Action:     "lte",
					Deprecated: "Lte",
				},
				{
					Name:       "AfterEquals",
					Action:     "gte",
					Deprecated: "Gte",
				},
			},
		},
	}
}
