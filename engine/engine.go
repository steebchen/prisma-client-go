package engine

import (
	"net/http"
	"os/exec"
	"time"
)

func NewEngine(schema string, hasBinaryTargets bool) *Engine {
	engine := &Engine{
		Schema:           schema,
		hasBinaryTargets: hasBinaryTargets,
	}

	engine.http = &http.Client{
		Timeout: 30 * time.Second,
	}

	return engine
}

type Engine struct {
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
}

// deprecated
func (e *Engine) ReplaceSchema(replace func(schema string) string) {
	e.Schema = replace(e.Schema)
}
