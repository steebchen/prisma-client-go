package generator

import (
	"github.com/prisma/prisma-client-go/generator/dmmf"
	"github.com/prisma/prisma-client-go/generator/types"
)

// Root describes the generator output root.
type Root struct {
	Generator       Generator     `json:"generator"`
	OtherGenerators []Generator   `json:"otherGenerators"`
	SchemaPath      string        `json:"schemaPath"`
	DMMF            dmmf.Document `json:"DMMF"`
	Datasources     []Datasource  `json:"datasources"`
	// Datamodel provides the raw string of the Prisma datamodel.
	Datamodel string `json:"datamodel"`
	// BinaryPaths (optional)
	BinaryPaths BinaryPaths `json:"binaryPaths"`
}

// Config describes the options for the Prisma Client Go generator
type Config struct {
	Package types.String `json:"package"`
}

// Generator describes a generator defined in the Prisma schema.
type Generator struct {
	// Output (optional) holds the file path of where the client gets generated in.
	Output        string       `json:"output"`
	Name          types.String `json:"name"`
	Provider      types.String `json:"provider"`
	Config        Config       `json:"config"`
	BinaryTargets []string     `json:"binaryTargets"`
	// PinnedBinaryTarget (optional)
	PinnedBinaryTarget string `json:"pinnedBinaryTarget"`
}

// ConnectorType describes the Database of this generator.
type ConnectorType string

// ConnectorType values
const (
	ConnectorTypeMySQL      ConnectorType = "mysql"
	ConnectorTypeMongo      ConnectorType = "mongo"
	ConnectorTypeSQLite     ConnectorType = "sqlite"
	ConnectorTypePostgreSQL ConnectorType = "postgresql"
)

// Datasource describes a Prisma data source of any database type.
type Datasource struct {
	Name          types.String  `json:"name"`
	ConnectorType ConnectorType `json:"connectorType"`
	URL           EnvValue      `json:"url"`
	Config        interface{}   `json:"config"`
}

// EnvValue contains a string value and optionally information if, and if yes from where, an env var is used for this value.
type EnvValue struct {
	// FromEnvVar (optional)
	FromEnvVar string `json:"fromEnvVar"`
	Value      string `json:"value"`
}

// BinaryPaths holds the information of the paths to the Prisma binaries.
type BinaryPaths struct {
	// MigrationEngine (optional)
	MigrationEngine map[string]string `json:"migrationEngine"` // key target, value path
	// QueryEngine (optional)
	QueryEngine map[string]string `json:"queryEngine"`
	// IntrospectionEngine (optional)
	IntrospectionEngine map[string]string `json:"introspectionEngine"`
}
