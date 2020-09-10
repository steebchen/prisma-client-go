// Package platform provides runtime methods to find out the correct prisma binary to use
package platform

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var binaryNameWithSSLCache string

// BinaryPlatformName returns the name of the prisma binary which should be used,
// for example "darwin" or "linux-openssl-1.1.x"
func BinaryPlatformName() string {
	if binaryNameWithSSLCache != "" {
		return binaryNameWithSSLCache
	}

	platform := Name()

	if platform != "linux" {
		return platform
	}

	distro := getLinuxDistro()

	if distro == "alpine" {
		return "linux-musl"
	}

	ssl := getOpenSSL()

	name := fmt.Sprintf("%s-openssl-%s", distro, ssl)

	binaryNameWithSSLCache = name

	return name
}

// Name returns the platform name
func Name() string {
	return runtime.GOOS
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
	r := regexp.MustCompile(`^OpenSSL\s(\d+\.\d+)\.\d+`)
	matches := r.FindStringSubmatch(str)
	if len(matches) > 0 {
		return matches[1] + ".x"
	}
	// default to 1.1.x
	return "1.1.x"
}
