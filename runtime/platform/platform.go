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

// GetMame returns the name of the prisma binary which should be used
func GetName() string {
	platform := runtime.GOOS

	if platform != "linux" {
		return platform
	}

	distro := getLinuxDistro()
	ssl := getOpenSSL()

	return fmt.Sprintf("%s-openssl-%s", distro, ssl)
}

func getLinuxDistro() string {
	out, _ := exec.Command("cat", "/etc/os-release").CombinedOutput()

	if out == nil {
		if true {
			log.Printf("`cat /etc/os-release` failed")
		}
		return ""
	}

	return parseLinuxDistro(string(out))
}

func parseLinuxDistro(str string) string {
	log.Printf("str %s", str)
	var id string
	var idLike string

	// match everything after ID= except quotes and newlines
	idMatches := regexp.MustCompile(`(?m)^ID="?([^"\n]*)"?`).FindStringSubmatch(str)
	log.Printf("idMatches %+v", idMatches)
	if len(idMatches) > 0 {
		id = idMatches[1]
	}

	// match everything after ID_LIKE= except quotes and newlines
	idLikeMatches := regexp.MustCompile(`(?m)^ID_LIKE="?([^"\n]*)"?`).FindStringSubmatch(str)
	log.Printf("idLikeMatches %+v", idLikeMatches)
	if len(idLikeMatches) > 0 {
		idLike = idLikeMatches[1]
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
		if true {
			log.Printf("`openssl version` failed")
		}
		return ""
	}

	return parseOpenSSLVersion(string(out))
}

func parseOpenSSLVersion(str string) string {
	r := regexp.MustCompile(`^OpenSSL\s(\d+\.\d+)\.\d+`)
	matches := r.FindStringSubmatch(str)
	if len(matches) > 0 {
		return matches[1]
	}
	// default to 1.1
	return "1.1"
}
