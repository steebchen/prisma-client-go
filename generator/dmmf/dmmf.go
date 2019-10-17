package dmmf

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

type Document struct {
	Datamodel Datamodel `json:"datamodel"`
	Schema    Schema    `json:"schema"`
	Mappings  []Mapping `json:"mappings"`
}

type Enum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
	// DBName (optional)
	DBName string `json:"dBName"` // can also be undefined
}

type Datamodel struct {
	Models []Model `json:"models"`
	Enums  []Enum  `json:"enums"`
}

type Model struct {
	Name       string `json:"name"`
	IsEmbedded bool   `json:"isEmbedded"`
	// DBName (optional)
	DBName string  `json:"dbName"` // can also be undefined
	Fields []Field `json:"fields"`
}

type Field struct {
	Kind       DatamodelFieldKind `json:"kind"`
	Name       string             `json:"name"`
	IsRequired bool               `json:"isRequired"`
	IsList     bool               `json:"isList"`
	IsUnique   bool               `json:"isUnique"`
	IsId       bool               `json:"isId"`
	Type       string             `json:"type"`
	// DBName (optional)
	DBName      string `json:"dBName"` // can also be undefined
	IsGenerated bool   `json:"isGenerated"`
	// RelationToFields (optional)
	RelationToFields []interface{} `json:"relationToFields"`
	// RelationOnDelete (optional)
	RelationOnDelete string
	// RelationName (optional)
	RelationName string
}

type Schema struct {
	// RootQueryType (optional)
	RootQueryType string `json:"rootQueryType"`
	// RootMutationType (optional)
	RootMutationType string       `json:"rootMutationType"`
	InputTypes       []InputType  `json:"inputTypes"`
	OutputTypes      []OutputType `json:"outputTypes"`
	Enums            []Enum       `json:"enums"`
}

type QueryOutput struct {
	Name       string `json:"name"`
	IsRequired bool   `json:"isRequired"`
	IsList     bool   `json:"isList"`
}

type SchemaArg struct {
	Name      string          `json:"name"`
	InputType SchemaInputType `json:"inputType"`
	// IsRelationFilter (optional)
	IsRelationFilter bool `json:"isRelationFilter"`
}

type SchemaInputType struct {
	IsRequired bool      `json:"isRequired"`
	IsList     bool      `json:"isList"`
	Type       string    `json:"type"` // this was declared as ArgType
	Kind       FieldKind `json:"kind"`
}

type OutputType struct {
	Name   string        `json:"name"`
	Fields []SchemaField `json:"fields"`
	// IsEmbedded (optional)
	IsEmbedded bool `json:"isEmbedded"`
}

type SchemaField struct {
	Name       string           `json:"name"`
	OutputType SchemaOutputType `json:"outputType"`
	Args       []SchemaArg      `json:"args"`
}

type SchemaOutputType struct {
	Type       string    `json:"type"` // note that in the serialized state we don't have the reference to MergedOutputTypes
	IsList     bool      `json:"isList"`
	IsRequired bool      `json:"isRequired"`
	Kind       FieldKind `json:"kind"`
}

type InputType struct {
	Name string `json:"name"`
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
	Model string `json:"model"`
	// FindOne (optional)
	FindOne string `json:"findOne"`
	// FindMany (optional)
	FindMany string `json:"findMany"`
	// Create (optional)
	Create string `json:"create"`
	// Update (optional)
	Update string `json:"update"`
	// UpdateMany (optional)
	UpdateMany string `json:"updateMany"`
	// Upsert (optional)
	Upsert string `json:"upsert"`
	// Delete (optional)
	Delete string `json:"delete"`
	// DeleteMany (optional)
	DeleteMany string `json:"deleteMany"`
}

type ModelAction string

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
