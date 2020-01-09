package engine

import (
	"net/http"
	"os/exec"
	"time"
)

func NewEngine(schema string) *Engine {
	e := &Engine{
		schema: schema,
	}

	e.http = &http.Client{
		Timeout: 30 * time.Second,
	}

	return e
}

type Engine struct {
	// cmd holds the prisma binary process
	cmd *exec.Cmd

	// http is the internal http client
	http *http.Client

	// url holds the query-engine url
	url string

	// schema contains the prisma schema
	schema string
}
