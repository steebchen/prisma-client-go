package engine

import (
	"net/http"
	"os/exec"
	"sync"
)

func NewQueryEngine(schema string, hasBinaryTargets bool, datasources string, datasourceURL string) *QueryEngine {
	return &QueryEngine{
		Schema:           schema,
		hasBinaryTargets: hasBinaryTargets,
		datasources:      datasources,
		datasourceURL:    datasourceURL,
		http:             &http.Client{},
	}
}

type QueryEngine struct {
	// Schema contains the prisma Schema
	Schema string

	// cmd holds the prisma binary process
	cmd *exec.Cmd

	// http is the internal http client
	http *http.Client

	// datasources holds the raw datasources
	datasources string

	// datasourceURL holds the sanitized datasourceURL which is overridden in the datasource above
	datasourceURL string

	// httpURL holds the query-engine httpURL
	httpURL string

	// hasBinaryTargets can be toggled by generated code from Schema.prisma whether binaryTargets
	// were specified and thus expects binaries in the local path
	hasBinaryTargets bool

	// connected indicates whether the user has called Connect()
	connected bool

	// disconnected indicates whether the user has called Disconnect()
	disconnected bool

	// closed keeps track of query engine status
	closed chan interface{}

	// onEngineError is a listener on the query engine stderr output
	onEngineError chan string

	// lastEngineError contains the last received error
	lastEngineError string

	mu sync.RWMutex
}

func (e *QueryEngine) Name() string {
	return "query-engine"
}

// deprecated
func (e *QueryEngine) ReplaceSchema(replace func(schema string) string) {
	e.Schema = replace(e.Schema)
}
