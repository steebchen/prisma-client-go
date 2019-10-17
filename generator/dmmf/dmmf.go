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
	Datamodel Datamodel
	Schema    Schema
	Mappings  []Mapping
}

type Enum struct {
	Name   string
	Values []string
	DBName *string // can also be undefined
}

type Datamodel struct {
	Models []Model
	Enums  []Enum
}

type Model struct {
	Name       string
	IsEmbedded bool
	DBName     *string // can also be undefined
	Fields     []Field
}

type Field struct {
	Kind             DatamodelFieldKind
	Name             string
	IsRequired       bool
	IsList           bool
	IsUnique         bool
	IsId             bool
	Type             string
	DBName           *string // can also be undefined
	IsGenerated      bool
	RelationToFields *[]interface{}
	RelationOnDelete *string
	RelationName     *string
}

type Schema struct {
	RootQueryType    *string
	RootMutationType *string
	InputTypes       []InputType
	OutputTypes      []OutputType
	Enums            []Enum
}

type QueryOutput struct {
	Name       string
	IsRequired bool
	IsList     bool
}

type SchemaArg struct {
	Name             string
	InputType        SchemaInputType
	IsRelationFilter *bool
}

type SchemaInputType struct {
	IsRequired bool
	IsList     bool
	Type       string // this was declared as ArgType
	Kind       FieldKind
}

type OutputType struct {
	Name       string
	Fields     []SchemaField
	IsEmbedded *bool
}

type SchemaField struct {
	Name       string
	OutputType SchemaOutputType
	Args       []SchemaArg
}

type SchemaOutputType struct {
	Type       string // note that in the serialized state we don't have the reference to MergedOutputTypes
	IsList     bool
	IsRequired bool
	Kind       FieldKind
}

type InputType struct {
	Name        string
	IsWhereType *bool // this is needed to transform it back
	IsOrderType *bool
	AtLeastOne  *bool
	AtMostOne   *bool
	Fields      []SchemaArg
}

type Mapping struct {
	Model      string
	FindOne    *string
	FindMany   *string
	Create     *string
	Update     *string
	UpdateMany *string
	Upsert     *string
	Delete     *string
	DeleteMany *string
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
