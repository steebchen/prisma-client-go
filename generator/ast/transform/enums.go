package transform

type Enum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func (r *AST) enums() []Enum {
	var enums []Enum
	for _, enum := range r.dmmf.Schema.EnumTypes.Enums {
		enums = append(enums, Enum{
			Name:   enum.Name,
			Values: enum.Values,
		})
	}
	return enums
}
