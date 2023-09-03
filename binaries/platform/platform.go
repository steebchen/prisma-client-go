// Package platform provides runtime methods to find out the correct prisma binary to use
package platform

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var binaryNameWithSSLCache string

// BinaryPlatformNameDynamic returns the name of the prisma binary which should be used,
// for example "darwin" or "linux-openssl-1.1.x". This can include dynamically linked binaries.
func BinaryPlatformNameDynamic() string {
	if binaryNameWithSSLCache != "" {
		return binaryNameWithSSLCache
	}

	platform := Name()
	arch := Arch()

	// other supported platforms are darwin and windows
	if platform != "linux" {
		// special case for darwin arm64
		if platform == "darwin" && arch == "arm64" {
			return "darwin-arm64"
		}
		// otherwise, return `darwin` or `windows`
		return platform
	}

	distro := getLinuxDistro()

	ssl := getOpenSSL()

	name := fmt.Sprintf("%s-openssl-%s", distro, ssl)

	binaryNameWithSSLCache = name

	return name
}

// BinaryPlatformNameStatic returns the name of the prisma binary which should be used,
// for example "darwin" or "linux-static-x64". This only includes statically linked binaries.
func BinaryPlatformNameStatic() string {
	platform := Name()
	arch := Arch()

	// other supported platforms are darwin and windows
	if platform != "linux" {
		// special case for darwin arm64
		if platform == "darwin" && arch == "arm64" {
			return "darwin-arm64"
		}
		// otherwise, return `darwin` or `windows`
		return platform
	}

	return fmt.Sprintf("linux-static-%s", arch)
}

// Name returns the platform name
func Name() string {
	return runtime.GOOS
}

func Arch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x64"
	case "arm64":
		return "arm64"
	default:
		log.Printf("warning: unsupported architecture %s, falling back to x64", runtime.GOARCH)
		return "x64"
	}
}

// CheckForExtension adds a .exe extension on windows (e.g. .gz -> .exe.gz)
func CheckForExtension(platform, path string) string {
	if platform == "windows" {
		if strings.Contains(path, ".gz") {
			return strings.Replace(path, ".gz", ".exe.gz", 1)
		}

		return path + ".exe"
	}

	return path
}

func getLinuxDistro() string {
	out, _ := exec.Command("cat", "/etc/os-release").CombinedOutput()

	if out != nil {
		return parseLinuxDistro(string(out))
	}

	return "debian"
}

func parseLinuxDistro(str string) string {
	var id string
	var idLike string

	// match everything after `ID=` except quotes and newlines
	idMatches := regexp.MustCompile(`(?m)^ID="?([^"\n]*)"?`).FindStringSubmatch(str)
	if len(idMatches) > 0 {
		id = idMatches[1]
	}

	// match everything after `ID_LIKE=` except quotes and newlines
	idLikeMatches := regexp.MustCompile(`(?m)^ID_LIKE="?([^"\n]*)"?`).FindStringSubmatch(str)
	if len(idLikeMatches) > 0 {
		idLike = idLikeMatches[1]
	}

	if id == "alpine" {
		return "alpine"
	}

	if strings.Contains(idLike, "centos") ||
		strings.Contains(idLike, "fedora") ||
		strings.Contains(idLike, "rhel") ||
		id == "fedora" {
		return "rhel"
	}

	if strings.Contains(idLike, "debian") ||
		strings.Contains(idLike, "ubuntu") ||
		id == "debian" {
		return "debian"
	}

	// default to debian as it's most common
	return "debian"
}

func getOpenSSL() string {
	out, _ := exec.Command("openssl", "version", "-v").CombinedOutput()

	if out == nil {
		return ""
	}

	return parseOpenSSLVersion(string(out))
}

// parseOpenSSLVersion returns the OpenSSL version, ignoring the patch version; e.g. 1.1.x
func parseOpenSSLVersion(str string) string {
	// parse minor and major version for OpenSSL 1.x
	r := regexp.MustCompile(`^OpenSSL\s(\d+\.\d+)\.\d+`)
	matches := r.FindStringSubmatch(str)
	if len(matches) > 1 && strings.HasPrefix(matches[1], "1.") {
		return matches[1] + ".x"
	}
	// parse major version for others
	r = regexp.MustCompile(`^OpenSSL\s(\d+)\.\d+\.\d+`)
	matches = r.FindStringSubmatch(str)
	if len(matches) > 0 {
		return matches[1] + ".0.x"
	}
	// default to 3.0.x
	return "3.0.x"
}
