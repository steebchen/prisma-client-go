package types

import (
	"github.com/iancoleman/strcase"
	"github.com/takuoki/gocase"
)

// String acts as a builtin string but provides useful casing methods
type String string

func (s String) GoCase() string {
	return gocase.To(strcase.ToCamel(string(s)))
}

func (s String) GoLowerCase() string {
	return gocase.To(strcase.ToLowerCamel(string(s)))
}

var builtin = map[string]string{
	"String":   "string",
	"Boolean":  "bool",
	"Int":      "int",
	"Float":    "float64",
	"DateTime": "time.Time",
}

// String acts as a builtin string but provides useful methods for type DMMF values
type Type string

func (t Type) Value() string {
	str := string(t)
	v, ok := builtin[str]
	if !ok {
		return str
	}
	return v
}
