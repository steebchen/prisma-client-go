package generator

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/steebchen/prisma-client-go/generator/ast/dmmf"
	"github.com/steebchen/prisma-client-go/generator/ast/transform"
	"github.com/steebchen/prisma-client-go/generator/types"
	"github.com/steebchen/prisma-client-go/logger"
)

// Root describes the generator output root.
type Root struct {
	Generator       Generator   `json:"generator"`
	OtherGenerators []Generator `json:"otherGenerators"`
	SchemaPath      string      `json:"schemaPath"`
	// Version contains the version hash of the Prisma engine
	Version     string        `json:"version"`
	DMMF        dmmf.Document `json:"DMMF"`
	Datasources []Datasource  `json:"datasources"`
	// Datamodel provides the raw string of the Prisma datamodel.
	Datamodel string `json:"datamodel"`
	// BinaryPaths (optional)
	BinaryPaths BinaryPaths    `json:"binaryPaths"`
	AST         *transform.AST `json:"ast"`
}

func (r *Root) EscapedDatamodel() string {
	return strings.ReplaceAll(r.Datamodel, "`", "'")
}

func (r *Root) GetDatasourcesJSON() string {
	ds := r.Datasources[0]

	data, err := json.Marshal([]Datasource{ds})
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (r *Root) GetEngineType() string {
	if str := os.Getenv("PRISMA_CLIENT_ENGINE_TYPE"); str != "" {
		return str
	}
	if str := r.Generator.Config.EngineType; str != "" {
		return str
	}
	return "binary"
}

// Config describes the options for the Prisma Client Go generator
type Config struct {
	EngineType        string       `json:"engineType"`
	Package           types.String `json:"package"`
	DisableGitignore  string       `json:"disableGitignore"`
	DisableGoBinaries string       `json:"disableGoBinaries"`
}

// Generator describes a generator defined in the Prisma schema.
type Generator struct {
	// Output holds the file path of where the client gets generated in.
	Output        *Value         `json:"output"`
	Name          types.String   `json:"name"`
	Provider      *Value         `json:"provider"`
	Config        Config         `json:"config"`
	BinaryTargets []BinaryTarget `json:"binaryTargets"`
	// PinnedBinaryTarget (optional)
	PinnedBinaryTarget string `json:"pinnedBinaryTarget"`
}

type BinaryTarget struct {
	FromEnvVar string `json:"fromEnvVar"`
	Value      string `json:"value"`
}

type Value struct {
	FromEnvVar string `json:"fromEnvVar"`
	Value      string `json:"value"`
}

// Provider describes the Database of this datasource.
type Provider string

// Provider values
//
//goland:noinspection GoUnusedConst
const (
	ProviderMySQL      Provider = "mysql"
	ProviderMongo      Provider = "mongo"
	ProviderSQLite     Provider = "sqlite"
	ProviderPostgreSQL Provider = "postgresql"
)

// Datasource describes a Prisma data source of any database type.
type Datasource struct {
	Name           types.String `json:"name"`
	Provider       Provider     `json:"provider"`
	ActiveProvider Provider     `json:"activeProvider"`
	URL            EnvValue     `json:"url"`
	Config         interface{}  `json:"config"`
}

// EnvValue contains a string value and optionally information if, and if yes from where, an env var is used for this value.
type EnvValue struct {
	// FromEnvVar (optional)
	FromEnvVar string `json:"fromEnvVar"`
	Value      string `json:"value"`
}

func (r *Root) GetSanitizedDatasourceURL() string {
	ds := r.Datasources[0]

	url := ds.URL.Value
	if ds.ActiveProvider != ProviderSQLite {
		return url
	}

	url = strings.ReplaceAll(url, "file:", "")
	url = strings.ReplaceAll(url, "sqlite:", "")

	if path.IsAbs(url) {
		return "file:" + url
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// get prisma schema path from prisma schema file
	schemaPath := path.Dir(r.SchemaPath)

	// trim /private as it is some kind of symlink on macOS
	schemaPath = strings.Replace(schemaPath, "/private", "", 1)

	// replace /schema.prisma as we need just the directory
	schemaPath = strings.Replace(schemaPath, "schema.prisma", "", 1)

	// use the schema path to locate the sqlite file (as the path is relative to the schema)
	url = path.Join(schemaPath, url)

	// replace absolute URL to relative
	url = strings.Replace(url, wd, "", 1)

	url = strings.Trim(url, "/")

	// prefix with sqlite: to make it a valid connection string again
	url = "file:" + url

	logger.Debug.Printf("sanitizing relative sqlite path %s\n", url)

	return url
}

// BinaryPaths holds the information of the paths to the Prisma binaries.
type BinaryPaths struct {
	// MigrationEngine (optional)
	MigrationEngine map[string]string `json:"migrationEngine"` // key target, value path
	// QueryEngine (optional)
	QueryEngine map[string]string `json:"queryEngine"`
}
