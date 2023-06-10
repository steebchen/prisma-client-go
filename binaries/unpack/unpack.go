package unpack

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/binaries/platform"
	"github.com/prisma/prisma-client-go/logger"
)

// TODO check checksum after expanding file

// noinspection GoUnusedExportedFunction
func Unpack(data []byte, name string, version string) {
	start := time.Now()

	filename := fmt.Sprintf("prisma-query-engine-%s", name)

	// TODO check if dev env/dev binary in ~/.prisma
	// TODO check if engine in local dir OR env var

	tempDir := binaries.GlobalUnpackDir(version)

	file := platform.CheckForExtension(platform.Name(), path.Join(tempDir, filename))

	if err := os.MkdirAll(tempDir, 0750); err != nil {
		panic(fmt.Errorf("mkdirall failed: %w", err))
	}

	if _, err := os.Stat(file); err == nil {
		logger.Debug.Printf("query engine exists, not unpacking. %s", time.Since(start))
		return
	}

	f, err := os.Create(file)
	if err != nil {
		panic(fmt.Errorf("generate open go file: %w", err))
	}

	if _, err := fmt.Fprintf(f, "%s", string(data)); err != nil {
		panic(err)
	}

	if err := f.Close(); err != nil {
		panic(err)
	}

	if err := os.Chmod(file, os.ModePerm); err != nil {
		panic(fmt.Errorf("could not chmod +x %s: %w", file, err))
	}

	logger.Debug.Printf("unpacked at %s in %s", file, time.Since(start))
}
