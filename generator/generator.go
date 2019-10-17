package generator

import (
	"github.com/prisma/photongo/generator/dmmf"
)

type GeneratorConfig struct {
	Output             *string
	Name               string
	Provider           string
	Config             []string // Dictionary<string> // help
	BinaryTargets      []string
	PinnedBinaryTarget *string
}

type ConnectorType string

const (
	ConnectorTypeMySQL      ConnectorType = "mysql"
	ConnectorTypeMongo      ConnectorType = "mongo"
	ConnectorTypeSQLite     ConnectorType = "sqlite"
	ConnectorTypePostgreSQL ConnectorType = "postgresql"
)

type Datasource struct {
	Name          string
	ConnectorType ConnectorType
	Url           string // EnvValue // help
	Config        interface{}
}

type GeneratorOptions struct {
	Generator       GeneratorConfig
	OtherGenerators []GeneratorConfig
	SchemaPath      string
	DMMF            dmmf.Document
	Datasources     []Datasource
	Datamodel       string
	BinaryPaths     *BinaryPaths
}

type BinaryPaths struct {
	MigrationEngine     *map[string]string // key target, value path
	QueryEngine         *map[string]string
	IntrospectionEngine *map[string]string
}
