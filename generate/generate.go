package generate

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/prisma/photongo/binaries"
	"github.com/prisma/photongo/logger"
)

// Run the generator
func Run() error {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logger.L.Printf("using dir %s", wd)

	if err := binaries.Fetch(wd); err != nil {
		return fmt.Errorf("could not fetch binaries: %w", err)
	}

	prisma := binaries.PrismaCLIName()

	logger.L.Printf("running %s %s", path.Join(wd, prisma), "generate")

	cmd := exec.Command(path.Join(wd, prisma), "generate")
	queryEngine := wd + "/prisma-query-engine"
	migrationEngine := wd + "/prisma-migration-engine"
	cmd.Env = append(
		os.Environ(),
		"PRISMA_QUERY_ENGINE_BINARY="+queryEngine,
		"PRISMA_MIGRATION_ENGINE_BINARY="+migrationEngine,
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
