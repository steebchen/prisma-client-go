package dmmf

import (
	"fmt"

	"github.com/prisma/photongo/generator/types"
)

type FieldKind string

const (
	FieldKindScalar FieldKind = "scalar"
	FieldKindObject FieldKind = "object"
	FieldKindEnum   FieldKind = "enum"
)

type DatamodelFieldKind string

const (
	DatamodelFieldKindScalar   DatamodelFieldKind = "scalar"
	DatamodelFieldKindRelation DatamodelFieldKind = "relation"
	DatamodelFieldKindEnum     DatamodelFieldKind = "enum"
)

func (v DatamodelFieldKind) IncludeInStruct() bool {
	return v == DatamodelFieldKindScalar || v == DatamodelFieldKindEnum
}

type Document struct {
	Datamodel Datamodel `json:"datamodel"`
	Schema    Schema    `json:"schema"`
	Mappings  []Mapping `json:"mappings"`
}

type Enum struct {
	Name   types.String   `json:"name"`
	Values []types.String `json:"values"`
	// DBName (optional)
	DBName types.String `json:"dBName"`
}

type Datamodel struct {
	Models []Model `json:"models"`
	Enums  []Enum  `json:"enums"`
}

type Model struct {
	Name       types.String `json:"name"`
	IsEmbedded bool         `json:"isEmbedded"`
	// DBName (optional)
	DBName types.String `json:"dbName"`
	Fields []Field      `json:"fields"`
}

type Field struct {
	Kind       DatamodelFieldKind `json:"kind"`
	Name       types.String       `json:"name"`
	IsRequired bool               `json:"isRequired"`
	IsList     bool               `json:"isList"`
	IsUnique   bool               `json:"isUnique"`
	IsId       bool               `json:"isId"`
	Type       types.Type         `json:"type"`
	// DBName (optional)
	DBName      types.String `json:"dBName"`
	IsGenerated bool         `json:"isGenerated"`
	// RelationToFields (optional)
	RelationToFields []interface{} `json:"relationToFields"`
	// RelationOnDelete (optional)
	RelationOnDelete types.String
	// RelationName (optional)
	RelationName types.String
}

func (f Field) Tag() string {
	return fmt.Sprintf("`json:\"%s\"`", f.Name.GoLowerCase())
}

type Schema struct {
	// RootQueryType (optional)
	RootQueryType types.String `json:"rootQueryType"`
	// RootMutationType (optional)
	RootMutationType types.String `json:"rootMutationType"`
	InputTypes       []InputType  `json:"inputTypes"`
	OutputTypes      []OutputType `json:"outputTypes"`
	Enums            []Enum       `json:"enums"`
}

type QueryOutput struct {
	Name       types.String `json:"name"`
	IsRequired bool         `json:"isRequired"`
	IsList     bool         `json:"isList"`
}

type SchemaArg struct {
	Name      types.String    `json:"name"`
	InputType SchemaInputType `json:"inputType"`
	// IsRelationFilter (optional)
	IsRelationFilter bool `json:"isRelationFilter"`
}

type SchemaInputType struct {
	IsRequired bool         `json:"isRequired"`
	IsList     bool         `json:"isList"`
	Type       types.String `json:"type"` // this was declared as ArgType
	Kind       FieldKind    `json:"kind"`
}

type OutputType struct {
	Name   types.String  `json:"name"`
	Fields []SchemaField `json:"fields"`
	// IsEmbedded (optional)
	IsEmbedded bool `json:"isEmbedded"`
}

type SchemaField struct {
	Name       types.String     `json:"name"`
	OutputType SchemaOutputType `json:"outputType"`
	Args       []SchemaArg      `json:"args"`
}

type SchemaOutputType struct {
	Type       types.String `json:"type"` // note that in the serialized state we don't have the reference to MergedOutputTypes
	IsList     bool         `json:"isList"`
	IsRequired bool         `json:"isRequired"`
	Kind       FieldKind    `json:"kind"`
}

type InputType struct {
	Name types.String `json:"name"`
	// IsWhereType (optional)
	IsWhereType bool `json:"isWhereType"` // this is needed to transform it back
	// IsOrderType (optional)
	IsOrderType bool `json:"isOrderType"`
	// AtLeastOne (optional)
	AtLeastOne bool `json:"atLeastOne"`
	// AtMostOne (optional)
	AtMostOne bool        `json:"atMostOne"`
	Fields    []SchemaArg `json:"fields"`
}

type Mapping struct {
	Model types.String `json:"model"`
	// FindOne (optional)
	FindOne types.String `json:"findOne"`
	// FindMany (optional)
	FindMany types.String `json:"findMany"`
	// Create (optional)
	Create types.String `json:"create"`
	// Update (optional)
	Update types.String `json:"update"`
	// UpdateMany (optional)
	UpdateMany types.String `json:"updateMany"`
	// Upsert (optional)
	Upsert types.String `json:"upsert"`
	// Delete (optional)
	Delete types.String `json:"delete"`
	// DeleteMany (optional)
	DeleteMany types.String `json:"deleteMany"`
}

type ModelAction types.String

const (
	ModelActionFindOne    ModelAction = "findOne"
	ModelActionFindMany   ModelAction = "findMany"
	ModelActionCreate     ModelAction = "create"
	ModelActionUpdate     ModelAction = "update"
	ModelActionUpdateMany ModelAction = "updateMany"
	ModelActionUpsert     ModelAction = "upsert"
	ModelActionDelete     ModelAction = "delete"
	ModelActionDeleteMany ModelAction = "deleteMany"
)
