package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/prisma/photongo/binaries"
	"github.com/prisma/photongo/binaries/platform"
	"github.com/prisma/photongo/logger"
)

func (e *Engine) Connect() error {
	logger.L.Printf("ensure query engine binary...")

	startEngine := time.Now()

	logger.L.Printf("connecting to engine...")

	binariesPath := binaries.GlobalPath()
	binaryName := platform.BinaryNameWithSSL()

	var file string

	name := "prisma-query-engine-"
	localPath := "./" + path.Join(name+binaryName)
	globalPath := path.Join(binariesPath, name+binaryName)

	logger.L.Printf("expecting query engine `%s`", name+binaryName)

	// TODO write tests for all cases
	// TODO just check for binaryTargets? if binaryTargets are provided, only use that mechanism

	// first, check if the query engine binary is being overridden by PRISMA_QUERY_ENGINE_BINARY

	prismaQueryEngineBinary := os.Getenv("PRISMA_QUERY_ENGINE_BINARY")
	if prismaQueryEngineBinary != "" {
		logger.L.Printf("PRISMA_QUERY_ENGINE_BINARY is defined, using %s", prismaQueryEngineBinary)
		file = prismaQueryEngineBinary
		// TODO check with os.Stat if the binary exists
	} else if _, err := os.Stat(localPath); err == nil {
		// check in the local working directory
		logger.L.Printf("query engine found in working directory")
		file = localPath
	} else if _, err := os.Stat(globalPath); err == nil {
		// check in the global cache directory
		logger.L.Printf("query engine found in global path")
		file = globalPath
	} else {
		logger.L.Printf("no query engine defined or found. fetching it at runtime...")

		start := time.Now()

		to := binaries.GlobalPath()

		qe, err := binaries.DownloadEngine("query-engine", to)
		if err != nil {
			return fmt.Errorf("could not fetch query engine: %w", err)
		}

		file = qe

		logger.L.Printf("using query engine at %s", qe)
		logger.L.Printf("ensure query engine took %s", time.Since(start))
	}

	logger.L.Printf("launching query engine at %s", file)

	port, err := getPort()
	if err != nil {
		return fmt.Errorf("get free port: %w", err)
	}

	logger.L.Printf("running query-engine on port %s", port)

	e.url = "http://localhost:" + port

	e.cmd = exec.Command(file, "-p", port)

	e.cmd.Stdout = os.Stdout
	e.cmd.Stderr = os.Stderr

	e.cmd.Env = append(
		os.Environ(),
		"PRISMA_DML="+e.schema,
		"RUST_LOG=error",
		"RUST_LOG_FORMAT=json",
	)

	// TODO fine tune this using log levels
	if logger.Enabled {
		e.cmd.Env = append(
			e.cmd.Env,
			"PRISMA_LOG_QUERIES=y",
			"RUST_LOG=info",
		)
	}

	err = e.cmd.Start()
	if err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	ctx := context.Background()

	// send a basic readiness healthcheck and retry if unsuccessful
	var connectErr error
	var gqlErrors []GQLError
	for i := 0; i < 100; i++ {
		body, err := e.Request(ctx, "GET", "/status", map[string]interface{}{})

		if err != nil {
			connectErr = err
			logger.L.Printf("could not connect; retrying...")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		var response GQLResponse

		err = json.Unmarshal(body, &response)
		if err != nil {
			connectErr = err
			logger.L.Printf("could not unmarshal response; retrying...")
			time.Sleep(50 * time.Millisecond)
			continue
		}

		if response.Errors != nil {
			gqlErrors = response.Errors
			logger.L.Printf("could not connect due to gql errors; retrying...")
			time.Sleep(50 * time.Millisecond)
			continue
		}

		connectErr = nil
		gqlErrors = nil
		break
	}

	if connectErr != nil {
		return fmt.Errorf("readiness query error: %w", connectErr)
	}

	if gqlErrors != nil {
		return fmt.Errorf("readiness gql errors: %+v", gqlErrors)
	}

	logger.L.Printf("connecting took %s", time.Since(startEngine))

	logger.L.Printf("connected.")
	return nil
}

func (e *Engine) Disconnect() error {
	logger.L.Printf("disconnecting...")

	if err := e.cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("send signal: %w", err)
	}

	if err := e.cmd.Wait(); err != nil {
		// TODO: is this a bug in the query-engine?
		if err.Error() != "signal: interrupt" {
			return fmt.Errorf("wait for process: %w", err)
		}
	}

	logger.L.Printf("disconnected.")
	return nil
}
