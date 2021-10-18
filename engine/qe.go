package engine

import (
	"net/http"
	"os/exec"
)

func New(schema string, hasBinaryTargets bool) *QueryEngine {
	return &QueryEngine{
		Schema:           schema,
		hasBinaryTargets: hasBinaryTargets,
		http:             &http.Client{},
	}
}

type QueryEngine struct {
	// cmd holds the prisma binary process
	cmd *exec.Cmd

	// http is the internal http client
	http *http.Client

	// url holds the query-engine url
	url string

	// Schema contains the prisma Schema
	Schema string

	// hasBinaryTargets can be toggled by generated code from Schema.prisma whether binaryTargets
	// were specified and thus expects binaries in the local path
	hasBinaryTargets bool

	disconnected bool
}

func (e *QueryEngine) Name() string {
	return "query-engine"
}

// deprecated
func (e *QueryEngine) ReplaceSchema(replace func(schema string) string) {
	e.Schema = replace(e.Schema)
}
