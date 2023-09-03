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

	"github.com/joho/godotenv"

	"github.com/steebchen/prisma-client-go/binaries"
	"github.com/steebchen/prisma-client-go/binaries/platform"
	"github.com/steebchen/prisma-client-go/binaries/unpack"
	"github.com/steebchen/prisma-client-go/logger"
)

func (e *QueryEngine) Connect() error {
	logger.Debug.Printf("ensure query engine binary...")

	_ = godotenv.Load(".env")
	_ = godotenv.Load("db/.env")
	_ = godotenv.Load("prisma/.env")

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

func (e *QueryEngine) Disconnect() error {
	e.disconnected = true
	logger.Debug.Printf("disconnecting...")

	if platform.Name() == "windows" {
		if err := e.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("kill process: %w", err)
		}
		return nil
	}

	if err := e.cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("send signal: %w", err)
	}

	if err := e.cmd.Wait(); err != nil {
		if err.Error() != "signal: interrupt" {
			return fmt.Errorf("wait for process: %w", err)
		}
	}

	logger.Debug.Printf("disconnected.")
	return nil
}

func (e *QueryEngine) ensure() (string, error) {
	ensureEngine := time.Now()

	unpackPath := binaries.GlobalUnpackDir(binaries.EngineVersion)
	cachePath := binaries.GlobalCacheDir()

	// check for darwin/windows/linux first
	binaryName := platform.CheckForExtension(platform.Name(), platform.BinaryPlatformNameStatic())
	exactBinaryName := platform.CheckForExtension(platform.Name(), platform.BinaryPlatformNameDynamic())

	var file string
	// forceVersion saves whether a version check should be done, which should be disabled
	// when providing a custom query engine value
	forceVersion := true

	name := "prisma-query-engine-"
	localStatic := path.Join("./", name+binaryName)
	localExact := path.Join("./", name+exactBinaryName)
	globalUnpackStatic := path.Join(unpackPath, name+binaryName)
	globalUnpackExact := path.Join(unpackPath, name+exactBinaryName)
	cacheStatic := path.Join(cachePath, binaries.EngineVersion, name+binaryName)
	cacheExact := path.Join(cachePath, binaries.EngineVersion, name+exactBinaryName)

	logger.Debug.Printf("checking for local query engine `%s` or `%s`", localStatic, localExact)
	logger.Debug.Printf("checking for global query engine `%s` or `%s`", globalUnpackStatic, globalUnpackExact)
	logger.Debug.Printf("checking for cached query engine `%s` or `%s`", cacheStatic, cacheExact)

	// TODO write tests for all cases

	// first, check if the query engine binary is being overridden by PRISMA_QUERY_ENGINE_BINARY
	prismaQueryEngineBinary := os.Getenv("PRISMA_QUERY_ENGINE_BINARY")
	if prismaQueryEngineBinary != "" {
		logger.Debug.Printf("PRISMA_QUERY_ENGINE_BINARY is defined, using %s", prismaQueryEngineBinary)

		if _, err := os.Stat(prismaQueryEngineBinary); err != nil {
			return "", fmt.Errorf("PRISMA_QUERY_ENGINE_BINARY was provided, but no query engine was found at %s", prismaQueryEngineBinary)
		}

		file = prismaQueryEngineBinary
		forceVersion = false
	} else {
		if qe := os.Getenv(unpack.FileEnv); qe != "" {
			logger.Debug.Printf("using unpacked file env %s %s", unpack.FileEnv, qe)

			if info, err := os.Stat(qe); err == nil {
				file = qe
				logger.Debug.Printf("exact query engine found in working directory: %s %+v", file, info)
			} else {
				return "", fmt.Errorf("prisma query engine was expected at %s via FileEnv but was not found", qe)
			}
		}

		if info, err := os.Stat(localExact); err == nil {
			file = localExact
			logger.Debug.Printf("exact query engine found in working directory: %s %+v", file, info)
		} else if info, err = os.Stat(localStatic); err == nil {
			file = localStatic
			logger.Debug.Printf("query engine found in working directory: %s %+v", file, info)
		} else if info, err = os.Stat(cacheExact); err == nil {
			file = cacheExact
			logger.Debug.Printf("query engine found in cache path: %s %+v", file, info)
		} else if info, err = os.Stat(cacheStatic); err == nil {
			file = cacheStatic
			logger.Debug.Printf("exact query engine found in cache path: %s %+v", file, info)
		} else if info, err = os.Stat(globalUnpackExact); err == nil {
			file = globalUnpackExact
			logger.Debug.Printf("query engine found in global path: %s %+v", file, info)
		} else if info, err = os.Stat(globalUnpackStatic); err == nil {
			file = globalUnpackStatic
			logger.Debug.Printf("exact query engine found in global path: %s %+v", file, info)
		}
	}

	if file == "" {
		// TODO log instructions on how to fix this problem
		return "", fmt.Errorf("no binary found")
	}

	startVersion := time.Now()
	out, err := exec.Command(file, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("version check failed: %w", err)
	}
	logger.Debug.Printf("version check took %s", time.Since(startVersion))

	if v := strings.TrimSpace(strings.Replace(string(out), "query-engine", "", 1)); binaries.EngineVersion != v {
		note := "Did you forget to run `go run github.com/steebchen/prisma-client-go generate`?"
		msg := fmt.Errorf("expected query engine version `%s` but got `%s`\n%s", binaries.EngineVersion, v, note)
		if forceVersion {
			return "", msg
		}

		logger.Info.Printf("%s, ignoring since custom query engine was provided", msg)
	}

	logger.Debug.Printf("using query engine at %s", file)
	logger.Debug.Printf("ensure query engine took %s", time.Since(ensureEngine))

	return file, nil
}

func (e *QueryEngine) spawn(file string) error {
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
		"PRISMA_DML="+e.Schema,
		"RUST_LOG=error",
		"RUST_LOG_FORMAT=json",
		"PRISMA_CLIENT_ENGINE_TYPE=binary",
		"PRISMA_ENGINE_PROTOCOL=graphql",
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

	if err := e.cmd.Start(); err != nil {
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

		if err := json.Unmarshal(body, &response); err != nil {
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
