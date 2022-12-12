package cli

import (
	"fmt"
	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/binaries/platform"
	"github.com/prisma/prisma-client-go/logger"
	"os"
	"os/exec"
	"path/filepath"
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

	logger.Debug.Printf("running %s %+v", filepath.ToSlash(filepath.Join(dir, prisma)), arguments)

	cmd := exec.Command(filepath.ToSlash(filepath.Join(dir, prisma)), arguments...) //nolint:gosec
	binaryName := platform.CheckForExtension(platform.Name(), platform.BinaryPlatformName())

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PRISMA_HIDE_UPDATE_MESSAGE=true")
	cmd.Env = append(cmd.Env, "PRISMA_CLI_QUERY_ENGINE_TYPE=binary")

	for _, engine := range binaries.Engines {
		var value string

		if env := os.Getenv(engine.Env); env != "" {
			logger.Debug.Printf("overriding %s to %s", engine.Name, env)
			value = env
		} else {
			value = filepath.ToSlash(filepath.Join(dir, binaries.EngineVersion, fmt.Sprintf("prisma-%s-%s", engine.Name, binaryName)))
		}

		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", engine.Env, value))
	}

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
