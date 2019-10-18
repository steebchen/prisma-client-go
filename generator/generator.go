package generator

import (
	"github.com/prisma/photongo/generator/dmmf"
)

// Root describes the generator output root
type Root struct {
	Generator       Generator     `json:"generator"`
	OtherGenerators []Generator   `json:"otherGenerators"`
	SchemaPath      string        `json:"schemaPath"`
	DMMF            dmmf.Document `json:"DMMF"`
	Datasources     []Datasource  `json:"datasources"`
	Datamodel       string        `json:"datamodel"`
	// BinaryPaths (optional)
	BinaryPaths BinaryPaths `json:"binaryPaths"`
}

// Config is the data structure of what you can define in your schema.prisma file
type Config struct {
	Package string `json:"package"`
}

type Generator struct {
	// Output (optional)
	Output        string   `json:"output"`
	Name          string   `json:"name"`
	Provider      string   `json:"provider"`
	Config        Config   `json:"config"`
	BinaryTargets []string `json:"binaryTargets"`
	// PinnedBinaryTarget (optional)
	PinnedBinaryTarget string `json:"pinnedBinaryTarget"`
}

type ConnectorType string

const (
	ConnectorTypeMySQL      ConnectorType = "mysql"
	ConnectorTypeMongo      ConnectorType = "mongo"
	ConnectorTypeSQLite     ConnectorType = "sqlite"
	ConnectorTypePostgreSQL ConnectorType = "postgresql"
)

type Datasource struct {
	Name          string        `json:"name"`
	ConnectorType ConnectorType `json:"connectorType"`
	Url           EnvValue      `json:"url"` // formerly EnvValue
	Config        interface{}   `json:"config"`
}

type EnvValue struct {
	// FromEnvVar (optional)
	FromEnvVar string `json:"fromEnvVar"`
	Value      string `json:"value"`
}

type BinaryPaths struct {
	// MigrationEngine (optional)
	MigrationEngine map[string]string `json:"migrationEngine"` // key target, value path
	// QueryEngine (optional)
	QueryEngine map[string]string `json:"queryEngine"`
	// IntrospectionEngine (optional)
	IntrospectionEngine map[string]string `json:"introspectionEngine"`
}
