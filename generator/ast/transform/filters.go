package transform

import (
	"github.com/steebchen/prisma-client-go/generator/types"
)

// Method defines the method for the virtual types method
type Method struct {
	// Name of the filter method ot use publicly, such as `Equals` or `Contains`
	Name string
	// Action of the filter for internal use in the query engine, e.g. `equals` or `contains`
	Action string
	// IsList defines whether the filter accepts a scalar or a slice of scalars
	IsList bool
	// Deprecated contains a description of what else to use, as this method will be removed in the future
	// If empty, this method is not deprecated
	Deprecated string
	// Type describes the type for this method. If empty, default to the parent scalar type.
	Type types.Type
}

// Filter defines the data struct for the virtual types method
type Filter struct {
	// Name of a filter, which can be a scala like `Int`, or a field name like `Age`
	Name string
	// Methods describe filter methods, such as `Equals`, `In` or `Contains`
	Methods []Method
}
