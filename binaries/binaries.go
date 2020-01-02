package binaries

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/prisma/photongo/binaries/platform"
	"github.com/prisma/photongo/logger"
)

// PrismaVersion is a hardcoded version of the Prisma CLI.
const PrismaVersion = "2.0.0-alpha.443"

// EngineVersion is a hardcoded version of the Prisma Engine.
// The versions can be found under https://github.com/prisma/prisma-engine/commits/master.
const EngineVersion = "2eb5a63ad82e15dc2c248a0ac84dc28cd35542d6"

var PrismaURL = "https://prisma-photongo.s3-eu-west-1.amazonaws.com/%s-%s-%s.gz"
var EngineURL = "https://prisma-builds.s3-eu-west-1.amazonaws.com/master/%s/%s/%s.gz"

// init overrides URLs if env variables are specific for debugging purposes and to
// be able to provide a fallback if the links above should go down
func init() {
	if prismaURL, ok := os.LookupEnv("PRISMA_CLI_URL"); ok {
		PrismaURL = prismaURL
	}
	if engineURL, ok := os.LookupEnv("PRISMA_ENGINE_URL"); ok {
		EngineURL = engineURL
	}
}

// PrismaCLIName returns the local file path of where the CLI lives
func PrismaCLIName() string {
	variation := platform.Name()
	return fmt.Sprintf("prisma-cli-%s", variation)
}

// GlobalPath returns the path of where the CLI lives
func GlobalPath() string {
	temp := os.TempDir()
	return path.Join(temp, "prisma", "photongo-prisma-binaries", PrismaVersion)
}

func fetch(toDir string, engine string, binary string) error {
	logger.L.Printf("checking %s...", engine)

	to := path.Join(toDir, fmt.Sprintf("prisma-%s-%s", engine, binary))

	urlName := engine
	// the query-engine binary to on S3 is "prisma"
	if engine == "query-engine" {
		urlName = "prisma"
	}
	url := fmt.Sprintf(EngineURL, EngineVersion, binary, urlName)

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		logger.L.Printf("%s is cached", to)
		return nil
	}

	logger.L.Printf("%s is missing, downloading...", engine)

	if err := download(url, to); err != nil {
		return fmt.Errorf("could not download %s to %s: %w", url, to, err)
	}

	logger.L.Printf("%s done", engine)

	return nil
}

func FetchBinary(toDir string, engineName string, binaryName string) error {
	return fetch(toDir, engineName, binaryName)
}

// FetchNative fetches the Prisma binaries needed for the generator to a given directory
func FetchNative(toDir string) error {
	if toDir == "" {
		return fmt.Errorf("toDir must be provided")
	}

	if !strings.HasPrefix(toDir, "/") {
		return fmt.Errorf("toDir must be absolute")
	}

	if err := os.MkdirAll(toDir, os.ModePerm); err != nil {
		return fmt.Errorf("could not run MkdirAll on path %s: %w", toDir, err)
	}

	if err := DownloadCLI(toDir); err != nil {
		return fmt.Errorf("could not download engines: %w", err)
	}

	engines := []string{
		"query-engine",
		"migration-engine",
		"introspection-engine",
	}

	for _, e := range engines {
		if _, err := DownloadEngine(e, toDir); err != nil {
			return fmt.Errorf("could not download engines: %w", err)
		}
	}

	return nil
}

func DownloadCLI(toDir string) error {
	cli := PrismaCLIName()
	to := path.Join(toDir, cli)
	url := fmt.Sprintf(PrismaURL, "prisma-cli", PrismaVersion, platform.Name())

	if _, err := os.Stat(to); os.IsNotExist(err) {
		logger.L.Printf("prisma cli doesn't exist, fetching...")

		if err := download(url, to); err != nil {
			return fmt.Errorf("could not download %s to %s: %w", url, to, err)
		}
	} else {
		logger.L.Printf("prisma cli is cached")
	}

	return nil
}

func DownloadEngine(name string, toDir string) (file string, err error) {
	if err := os.MkdirAll(toDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not run MkdirAll on path %s: %w", toDir, err)
	}

	binaryName := platform.BinaryNameWithSSL()

	logger.L.Printf("checking %s...", name)

	to := path.Join(toDir, fmt.Sprintf("prisma-%s-%s", name, binaryName))

	urlName := name
	// the query-engine binary to on S3 is "prisma"
	if name == "query-engine" {
		urlName = "prisma"
	}
	url := fmt.Sprintf(EngineURL, EngineVersion, binaryName, urlName)

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		logger.L.Printf("%s is cached", to)
		return to, nil
	}

	logger.L.Printf("%s is missing, downloading...", name)

	startDownload := time.Now()
	if err := download(url, to); err != nil {
		return "", fmt.Errorf("could not download %s to %s: %w", url, to, err)
	}

	logger.L.Printf("download() took %s", time.Since(startDownload))

	logger.L.Printf("%s done", name)

	return to, nil
}

func download(url string, to string) error {
	// copy to temp file first
	dest := to + ".tmp"

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not get %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received code %d from %s: %+v", resp.StatusCode, url, string(out))
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", dest, err)
	}
	defer out.Close()

	if err := os.Chmod(dest, 0777); err != nil {
		return fmt.Errorf("could not chmod +x %s: %w", url, err)
	}

	g, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("could not create gzip reader: %w", err)
	}
	defer g.Close()

	if _, err := io.Copy(out, g); err != nil {
		return fmt.Errorf("could not copy %s: %w", url, err)
	}

	// temp file is ready, now copy to the original destination
	if err := copyFile(dest, to); err != nil {
		return fmt.Errorf("copy temp file: %w", err)
	}

	return nil
}

func copyFile(from string, to string) error {
	input, err := ioutil.ReadFile(from)
	if err != nil {
		return fmt.Errorf("readfile: %w", err)
	}

	err = ioutil.WriteFile(to, input, 0777)
	if err != nil {
		return fmt.Errorf("writefile: %w", err)
	}

	return nil
}
