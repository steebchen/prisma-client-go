package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/binaries/platform"
	"github.com/prisma/prisma-client-go/logger"
)

func (e *Engine) Connect() error {
	logger.Debug.Printf("ensure query engine binary...")

	startEngine := time.Now()

	file, err := e.ensure()
	if err != nil {
		return fmt.Errorf("ensure: %w", err)
	}

	if err := e.spawn(file); err != nil {
		return fmt.Errorf("spawn: %w", err)
	}

	logger.Debug.Printf("connecting took %s", time.Since(startEngine))
	logger.Debug.Printf("connected.")

	return nil
}

func (e *Engine) Disconnect() error {
	logger.Debug.Printf("disconnecting...")

	if platform.Name() != "windows" {
		if err := e.cmd.Process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("send signal: %w", err)
		}

		if err := e.cmd.Wait(); err != nil {
			// TODO: is this a bug in the query-engine?
			if err.Error() != "signal: interrupt" {
				return fmt.Errorf("wait for process: %w", err)
			}
		}
	}

	if platform.Name() == "windows" {

		if err := e.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("kill process: %w", err)
		}
	}

	logger.Debug.Printf("disconnected.")
	return nil
}

func (e *Engine) ensure() (string, error) {
	ensureEngine := time.Now()

	binariesPath := binaries.GlobalTempDir()
	binaryName := platform.CheckForExtension(platform.BinaryPlatformName())

	var file string
	// forceVersion saves whether a version check should be done, which should be disabled
	// when providing a custom query engine value
	forceVersion := true

	name := "prisma-query-engine-"
	localPath := "./" + path.Join(name+binaryName)
	globalPath := path.Join(binariesPath, name+binaryName)

	logger.Debug.Printf("expecting query engine `%s`", name+binaryName)

	// TODO write tests for all cases

	// first, check if the query engine binary is being overridden by PRISMA_QUERY_ENGINE_BINARY
	prismaQueryEngineBinary := os.Getenv("PRISMA_QUERY_ENGINE_BINARY")
	if prismaQueryEngineBinary != "" {
		logger.Debug.Printf("PRISMA_QUERY_ENGINE_BINARY is defined, using %s", prismaQueryEngineBinary)

		if _, err := os.Stat(prismaQueryEngineBinary); err == nil {
			return "", fmt.Errorf("PRISMA_QUERY_ENGINE_BINARY was provided, but no query engine was found at %s", prismaQueryEngineBinary)
		}

		file = prismaQueryEngineBinary
		forceVersion = false
	}

	if _, err := os.Stat(localPath); err == nil {
		// check in the local working directory
		logger.Debug.Printf("query engine found in working directory")
		file = localPath
	}

	if e.hasBinaryTargets && file == "" {
		logger.Debug.Printf("binaryTargets provided, but no query engine found at `%s`", name)
		return "", fmt.Errorf("binary targets were provided, but no query engine was found, please provide/upload the query engine `%s` in the project dir", name)
	}

	if _, err := os.Stat(globalPath); err == nil {
		// check in the global cache directory
		logger.Debug.Printf("query engine found in global path")
		file = globalPath
	}

	if file == "" {
		logger.Info.Printf("no query engine defined or found")
		logger.Info.Printf("if you want to pre-fetch the query engine for better startup performance, specify `binaryTargets = [\"native\"]` in your schema.prisma file under \"generator\" and upload the query engine with your application.")
		logger.Info.Printf("fetching the query engine now...")

		qe, err := binaries.DownloadEngine("query-engine", binariesPath)
		if err != nil {
			return "", fmt.Errorf("could not fetch query engine: %w", err)
		}

		logger.Info.Printf("done.")

		file = qe
	}

	startVersion := time.Now()
	out, err := exec.Command(file, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("version check failed: %w", err)
	}
	logger.Debug.Printf("version check took %s", time.Since(startVersion))

	if v := strings.TrimSpace(strings.Replace(string(out), "query-engine", "", 1)); binaries.EngineVersion != v {
		msg := fmt.Errorf("expected query engine version `%s` but got `%s`", binaries.EngineVersion, v)
		if forceVersion {
			return "", msg
		}

		logger.Info.Printf("%s, ignoring since custom query engine was provided", msg)
	}

	logger.Debug.Printf("using query engine at %s", file)
	logger.Debug.Printf("ensure query engine took %s", time.Since(ensureEngine))

	return file, nil
}

func (e *Engine) spawn(file string) error {
	port, err := getPort()
	if err != nil {
		return fmt.Errorf("get free port: %w", err)
	}

	logger.Debug.Printf("running query-engine on port %s", port)

	e.url = "http://localhost:" + port

	e.cmd = exec.Command(file, "-p", port, "--enable-raw-queries")

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

	logger.Debug.Printf("starting engine...")

	err = e.cmd.Start()
	if err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	logger.Debug.Printf("connecting to engine...")

	ctx := context.Background()

	// send a basic readiness healthcheck and retry if unsuccessful
	var connectErr error
	var gqlErrors []GQLError
	for i := 0; i < 100; i++ {
		body, err := e.Request(ctx, "GET", "/status", map[string]interface{}{})

		if err != nil {
			connectErr = err
			logger.Debug.Printf("could not connect; retrying...")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		var response GQLResponse

		err = json.Unmarshal(body, &response)
		if err != nil {
			connectErr = err
			logger.Debug.Printf("could not unmarshal response; retrying...")
			time.Sleep(50 * time.Millisecond)
			continue
		}

		if response.Errors != nil {
			gqlErrors = response.Errors
			logger.Debug.Printf("could not connect due to gql errors; retrying...")
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

	return nil
}
