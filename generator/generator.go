package generator

import (
	"github.com/prisma/photongo/generator/dmmf"
)

type Config struct {
	// Output (optional)
	Output        string   `json:"output"`
	Name          string   `json:"name"`
	Provider      string   `json:"provider"`
	Config        []string `json:"config"` // formerly Dictionary<string>
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
	Url           string        `json:"url"` // formerly EnvValue
	Config        interface{}   `json:"config"`
}

type Options struct {
	Generator       Config        `json:"generator"`
	OtherGenerators []Config      `json:"otherGenerators"`
	SchemaPath      string        `json:"schemaPath"`
	DMMF            dmmf.Document `json:"DMMF"`
	Datasources     []Datasource  `json:"datasources"`
	Datamodel       string        `json:"datamodel"`
	// BinaryPaths (optional)
	BinaryPaths BinaryPaths `json:"binaryPaths"`
}

type BinaryPaths struct {
	// MigrationEngine (optional)
	MigrationEngine map[string]string `json:"migrationEngine"` // key target, value path
	// QueryEngine (optional)
	QueryEngine map[string]string `json:"queryEngine"`
	// IntrospectionEngine (optional)
	IntrospectionEngine map[string]string `json:"introspectionEngine"`
}
