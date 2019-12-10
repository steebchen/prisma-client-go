package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/prisma/photongo/binaries"
	"github.com/prisma/photongo/binaries/platform"
	"github.com/prisma/photongo/logger"
)

// Run the prisma CLI with given arguments
func Run(arguments []string) error {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := binaries.Fetch(wd); err != nil {
		return fmt.Errorf("could not fetch binaries: %w", err)
	}

	prisma := binaries.PrismaCLIName()

	logger.L.Printf("running %s %+v", path.Join(wd, prisma), arguments)

	cmd := exec.Command(path.Join(wd, prisma), arguments...)
	binaryName := platform.BinaryNameWithSSL()
	queryEngine := wd + "/prisma-query-engine-" + binaryName
	migrationEngine := wd + "/prisma-migration-engine-" + binaryName
	introspectionEngine := wd + "/prisma-introspection-engine-" + binaryName
	cmd.Env = append(
		os.Environ(),
		"PRISMA_QUERY_ENGINE_BINARY="+queryEngine,
		"PRISMA_MIGRATION_ENGINE_BINARY="+migrationEngine,
		"PRISMA_INTROSPECTION_ENGINE_BINARY="+introspectionEngine,
	)
	if logger.Debug {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not run prisma generate: %w", err)
	}

	return nil
}
