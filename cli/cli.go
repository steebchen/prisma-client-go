package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/prisma/photongo/binaries"
	"github.com/prisma/photongo/binaries/platform"
	"github.com/prisma/photongo/logger"
)

// Run the prisma CLI with given arguments
func Run(arguments []string, output bool) error {
	logger.L.Printf("running cli with args %+v", arguments)

	dir := binaries.GlobalPath()

	if err := binaries.Fetch(dir); err != nil {
		return fmt.Errorf("could not fetch binaries: %w", err)
	}

	prisma := binaries.PrismaCLIName()

	logger.L.Printf("running %s %+v", path.Join(dir, prisma), arguments)

	cmd := exec.Command(path.Join(dir, prisma), arguments...)
	binaryName := platform.BinaryNameWithSSL()
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
