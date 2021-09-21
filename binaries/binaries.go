package binaries

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/prisma/prisma-client-go/binaries/platform"
	"github.com/prisma/prisma-client-go/logger"
)

// PrismaVersion is a hardcoded version of the Prisma CLI.
const PrismaVersion = "3.1.1"

// EngineVersion is a hardcoded version of the Prisma Engine.
// The versions can be found under https://github.com/prisma/prisma-engines/commits/master
const EngineVersion = "c22652b7e418506fab23052d569b85d3aec4883f"

// PrismaURL points to an S3 bucket URL where the CLI binaries are stored.
var PrismaURL = "https://prisma-photongo.s3-eu-west-1.amazonaws.com/%s-%s-%s.gz"

// EngineURL points to an S3 bucket URL where the Prisma engines are stored.
var EngineURL = "https://binaries.prisma.sh/all_commits/%s/%s/%s.gz"

type Engine struct {
	Name string
	Env  string
}

var Engines = []Engine{{
	"query-engine",
	"PRISMA_QUERY_ENGINE_BINARY",
}, {
	"migration-engine",
	"PRISMA_MIGRATION_ENGINE_BINARY",
}, {
	"introspection-engine",
	"PRISMA_INTROSPECTION_ENGINE_BINARY",
}, {
	"prisma-fmt",
	"PRISMA_FMT_BINARY",
}}

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

var baseDirName = path.Join("prisma", "binaries")

// GlobalTempDir returns the path of where the engines live
// internally, this is the global temp dir
func GlobalTempDir() string {
	temp := os.TempDir()
	logger.Debug.Printf("temp dir: %s", temp)

	return path.Join(temp, baseDirName, "engines", EngineVersion)
}

func GlobalUnpackDir() string {
	return path.Join(GlobalTempDir(), "unpacked")
}

// GlobalCacheDir returns the path of where the CLI lives
// internally, this is the global temp dir
func GlobalCacheDir() string {
	cache, err := os.UserCacheDir()
	if err != nil {
		panic(fmt.Errorf("could not read user cache dir: %w", err))
	}

	logger.Debug.Printf("global cache dir: %s", cache)

	return path.Join(cache, baseDirName, "cli", PrismaVersion)
}

func FetchEngine(toDir string, engineName string, binaryPlatformName string) error {
	logger.Debug.Printf("checking %s...", engineName)

	to := platform.CheckForExtension(binaryPlatformName, path.Join(toDir, EngineVersion, fmt.Sprintf("prisma-%s-%s", engineName, binaryPlatformName)))

	binaryPlatformRemoteName := binaryPlatformName
	if binaryPlatformRemoteName == "linux" {
		binaryPlatformRemoteName = "linux-musl"
	}
	url := platform.CheckForExtension(binaryPlatformName, fmt.Sprintf(EngineURL, EngineVersion, binaryPlatformRemoteName, engineName))

	logger.Debug.Printf("download url %s", url)

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		logger.Debug.Printf("%s is cached", to)
		return nil
	}

	logger.Debug.Printf("%s is missing, downloading...", engineName)

	if err := download(url, to); err != nil {
		return fmt.Errorf("could not download %s to %s: %w", url, to, err)
	}

	logger.Debug.Printf("%s done", engineName)

	return nil
}

// FetchNative fetches the Prisma binaries needed for the generator to a given directory
func FetchNative(toDir string) error {
	if toDir == "" {
		return fmt.Errorf("toDir must be provided")
	}

	if !filepath.IsAbs(toDir) {
		return fmt.Errorf("toDir must be absolute")
	}

	if err := DownloadCLI(toDir); err != nil {
		return fmt.Errorf("could not download engines: %w", err)
	}

	for _, e := range Engines {
		if _, err := DownloadEngine(e.Name, toDir); err != nil {
			return fmt.Errorf("could not download engines: %w", err)
		}
	}

	return nil
}

func DownloadCLI(toDir string) error {
	cli := PrismaCLIName()
	to := platform.CheckForExtension(platform.Name(), path.Join(toDir, cli))
	url := platform.CheckForExtension(platform.Name(), fmt.Sprintf(PrismaURL, "prisma-cli", PrismaVersion, platform.Name()))

	logger.Debug.Printf("ensuring CLI %s from %s to %s", cli, to, url)

	if _, err := os.Stat(to); os.IsNotExist(err) {
		logger.Info.Printf("prisma cli doesn't exist, fetching... (this might take a few minutes)")

		if err := download(url, to); err != nil {
			return fmt.Errorf("could not download %s to %s: %w", url, to, err)
		}

		logger.Info.Printf("prisma cli fetched successfully.")
	} else {
		logger.Debug.Printf("prisma cli is cached")
	}

	return nil
}

func GetEnginePath(dir, engine, binaryName string) string {
	return platform.CheckForExtension(binaryName, path.Join(dir, EngineVersion, fmt.Sprintf("prisma-%s-%s", engine, binaryName)))
}

func DownloadEngine(name string, toDir string) (file string, err error) {
	binaryName := platform.BinaryPlatformName()

	logger.Debug.Printf("checking %s...", name)

	to := platform.CheckForExtension(binaryName, path.Join(toDir, EngineVersion, fmt.Sprintf("prisma-%s-%s", name, binaryName)))

	url := platform.CheckForExtension(binaryName, fmt.Sprintf(EngineURL, EngineVersion, binaryName, name))

	logger.Debug.Printf("download url %s", url)

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		logger.Debug.Printf("%s is cached", to)
		return to, nil
	}

	logger.Debug.Printf("%s is missing, downloading...", name)

	startDownload := time.Now()
	if err := download(url, to); err != nil {
		return "", fmt.Errorf("could not download %s to %s: %w", url, to, err)
	}

	logger.Debug.Printf("download() took %s", time.Since(startDownload))

	logger.Debug.Printf("%s done", name)

	return to, nil
}

func download(url string, to string) error {
	if err := os.MkdirAll(path.Dir(to), os.ModePerm); err != nil {
		return fmt.Errorf("could not run MkdirAll on path %s: %w", to, err)
	}

	// copy to temp file first
	dest := to + ".tmp"

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return fmt.Errorf("could not get %s: %w", url, err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received code %d from %s: %+v", resp.StatusCode, url, string(out))
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", dest, err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer out.Close()

	if err := os.Chmod(dest, os.ModePerm); err != nil {
		return fmt.Errorf("could not chmod +x %s: %w", url, err)
	}

	g, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("could not create gzip reader: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer g.Close()

	if _, err := io.Copy(out, g); err != nil { //nolint:gosec
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

	if err := ioutil.WriteFile(to, input, os.ModePerm); err != nil {
		return fmt.Errorf("writefile: %w", err)
	}

	return nil
}
