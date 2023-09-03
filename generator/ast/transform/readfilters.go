package transform

import (
	"strings"

	"github.com/steebchen/prisma-client-go/generator/ast/dmmf"
	"github.com/steebchen/prisma-client-go/generator/types"
)

const list = "List"

func (r *AST) readFilters() []Filter {
	var filters []Filter
	for _, scalar := range r.Scalars {
		combinations := [][]string{
			{
				scalar + "ListFilter",
				scalar + "NullableListFilter",
			},
			{
				scalar + "Filter",
				scalar + "NullableFilter",
			},
		}

		for _, c := range combinations {
			p := r.pick(c...)
			if p == nil {
				continue
			}
			var fields []Method
			for _, field := range p.Fields {
				if method := convertField(field); method != nil {
					fields = append(fields, *method)
				}
			}
			s := scalar
			if strings.Contains(p.Name.String(), "ListFilter") {
				s += list
			}
			filters = append(filters, Filter{
				Name:    s,
				Methods: fields,
			})
		}
	}
	for _, enum := range r.Enums {
		p := r.pick(
			"Enum"+enum.Name.String()+"Filter",
			"Enum"+enum.Name.String()+"NullableFilter",
		)
		if p == nil {
			continue
		}

		var fields []Method
		for _, field := range p.Fields {
			if method := convertField(field); method != nil {
				fields = append(fields, *method)
			}
		}

		filters = append(filters, Filter{
			Name:    enum.Name.String(),
			Methods: fields,
		})
	}

	// order by relevance

	for i, m := range r.Models {
		p := r.pick(m.Name.String() + "OrderByRelevanceInput")
		if p == nil {
			continue
		}

		var methods []Method
		for _, field := range p.Fields {
			if method := convertField(field); method != nil {
				methods = append(methods, *method)
			}
		}

		filters = append(filters, Filter{
			Name:    p.Name.String(),
			Methods: methods,
		})

		// add pseudo model field Relevance_, so one can do
		// db.User.Relevance_.X
		r.Models[i].Fields = append(r.Models[i].Fields, Field{
			Prisma: true,
			Field: dmmf.Field{
				Name: "relevance",
				Kind: "scalar",
				Type: types.Type(p.Name.GoCase()),
			},
		})
	}

	return filters
}

// ReadFilter returns a filter for a read operation by scalar
func (r *AST) ReadFilter(scalar string, isList bool) *Filter {
	scalar = strings.Replace(scalar, "NullableFilter", "", 1)
	scalar = strings.Replace(scalar, "ReadFilter", "", 1)
	if isList {
		scalar += list
	}
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

	var typeName types.Type
	var isList bool
	for _, inputType := range field.InputTypes {
		if (inputType.Location == "scalar" || inputType.Location == "enumTypes") && inputType.Type != "Null" {
			typeName = inputType.Type
			if inputType.IsList {
				isList = true
			}
		}
	}
	return &Method{
		Name:   field.Name.GoCase(),
		Action: field.Name.String(),
		Type:   typeName,
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
