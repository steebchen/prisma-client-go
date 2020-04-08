package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/binaries/platform"
	"github.com/prisma/prisma-client-go/logger"
)

// Run the prisma CLI with given arguments
func Run(arguments []string, output bool) error {
	logger.Debug.Printf("running cli with args %+v", arguments)
	// TODO respect initial PRISMA_<name>_BINARY env
	// TODO optionally override CLI filepath using PRISMA_CLI_PATH

	dir := binaries.GlobalCacheDir()

	if err := binaries.FetchNative(dir); err != nil {
		return fmt.Errorf("could not fetch binaries: %w", err)
	}

	prisma := binaries.PrismaCLIName()

	logger.Debug.Printf("running %s %+v", path.Join(dir, prisma), arguments)

	cmd := exec.Command(path.Join(dir, prisma), arguments...)
	binaryName := platform.BinaryPlatformName()
	queryEngine := dir + "/prisma-query-engine-" + binaryName
	migrationEngine := dir + "/prisma-migration-engine-" + binaryName
	introspectionEngine := dir + "/prisma-introspection-engine-" + binaryName
	cmd.Env = append(
		os.Environ(),
		"PRISMA_QUERY_ENGINE_BINARY="+queryEngine,
		"PRISMA_MIGRATION_ENGINE_BINARY="+migrationEngine,
		"PRISMA_INTROSPECTION_ENGINE_BINARY="+introspectionEngine,
	)

	cmd.Stdin = os.Stdin

	if output {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not run %+v: %w", arguments, err)
	}

	return nil
}
