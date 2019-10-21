package types

import (
	"github.com/iancoleman/strcase"
	"github.com/takuoki/gocase"
)

// String acts as a builtin string but provides useful casing methods.
type String string

// GoCase transforms strings into Go-style casing, meaning uppercase including Go casing edge cases.
func (s String) GoCase() string {
	return gocase.To(strcase.ToCamel(string(s)))
}

// GoLowerCase transforms strings into Go-style lowercase casing. It's like GoCase but used for private fields.
func (s String) GoLowerCase() string {
	return gocase.To(strcase.ToLowerCamel(string(s)))
}

// builtin Go types
var builtin = map[string]string{
	"String":   "string",
	"Boolean":  "bool",
	"Int":      "int",
	"Float":    "float64",
	"DateTime": "time.Time",
}

// Type acts as a builtin string but provides useful methods for type DMMF values.
type Type string

// Value returns the native value of a type.
func (t Type) Value() string {
	str := string(t)
	v, ok := builtin[str]
	if !ok {
		return str
	}
	return v
}
