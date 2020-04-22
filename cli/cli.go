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

const QE = "PRISMA_QUERY_ENGINE_BINARY"
const ME = "PRISMA_MIGRATION_ENGINE_BINARY"
const IE = "PRISMA_INTROSPECTION_ENGINE_BINARY"

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
	binaryName := platform.CheckForExtension(platform.BinaryPlatformName())
	queryEngine := dir + "/prisma-query-engine-" + binaryName
	migrationEngine := dir + "/prisma-migration-engine-" + binaryName
	introspectionEngine := dir + "/prisma-introspection-engine-" + binaryName

	if qe := os.Getenv(QE); qe != "" {
		logger.Debug.Printf("overriding query engine to %s", qe)
		queryEngine = qe
	}

	if me := os.Getenv(ME); me != "" {
		logger.Debug.Printf("overriding migration engine to %s", me)
		migrationEngine = me
	}

	if ie := os.Getenv(IE); ie != "" {
		logger.Debug.Printf("overriding introspection engine to %s", ie)
		introspectionEngine = ie
	}

	cmd.Env = append(
		os.Environ(),
		QE+"="+queryEngine,
		ME+"="+migrationEngine,
		IE+"="+introspectionEngine,
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
