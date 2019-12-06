package binaries

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/prisma/photongo/binaries/platform"
	"github.com/prisma/photongo/logger"
)

// PrismaVersion is a hardcoded version of the Prisma CLI.
const PrismaVersion = "1"

// EngineVersion is a hardcoded version of the Prisma Engine.
// The versions can be found under https://github.com/prisma/prisma-engine/commits/master.
const EngineVersion = "4028eec09329a14692b13f06581329fddb7b2876"

const PrismaURL = "https://prisma-binaries-photongo.s3.eu-central-1.amazonaws.com/%s.gz"
const EngineURL = "https://prisma-builds.s3-eu-west-1.amazonaws.com/master/%s/%s/%s.gz"

// PrismaCLIName returns the local file path of where the CLI is located
func PrismaCLIName() string {
	variation := platform.Name()
	return fmt.Sprintf("prisma-cli-%s-%s", variation, PrismaVersion)
}

// Fetch fetches the Prisma binaries needed for the generator to a given directory
func Fetch(toDir string) error {
	if toDir == "" {
		return fmt.Errorf("toDir must be provided")
	}

	if !strings.HasPrefix(toDir, "/") {
		return fmt.Errorf("toDir must be absolute")
	}

	cli := PrismaCLIName()
	to := path.Join(toDir, cli)
	url := fmt.Sprintf(PrismaURL, cli)
	if err := download(url, to); err != nil {
		return fmt.Errorf("could not download %s to %s: %w", url, to, err)
	}

	engines := []string{
		"query-engine",
		"migration-engine",
		"introspection-engine",
	}

	for _, e := range engines {
		to := path.Join(toDir, fmt.Sprintf("prisma-%s-%s", e, EngineVersion[:7]))

		urlName := e
		// the query-engine binary to on S3 is "prisma"
		if e == "query-engine" {
			urlName = "prisma"
		}
		url := fmt.Sprintf(EngineURL, EngineVersion, platform.BinaryNameWithSSL(), urlName)

		if err := download(url, to); err != nil {
			return fmt.Errorf("could not download %s to %s: %w", url, to, err)
		}
	}

	return nil
}

func download(url string, dest string) error {
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		logger.L.Printf("%s exists", dest)
		return nil
	}

	logger.L.Printf("downloading %s...", url)

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", dest, err)
	}
	defer out.Close()

	err = os.Chmod(dest, 0777)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not get %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received code %d from %s: %+v", resp.StatusCode, url, string(out))
	}

	g, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("could not create gzip reader: %w", err)
	}
	defer g.Close()

	if _, err := io.Copy(out, g); err != nil {
		return fmt.Errorf("could not copy %s: %w", url, err)
	}

	// verify that the binary is working
	cmd := exec.Command(dest, "--help")

	if logger.Debug {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not run %s: %w", dest, err)
	}

	logger.L.Printf("done")

	return nil
}
