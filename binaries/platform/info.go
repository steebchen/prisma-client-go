package platform

import "strings"

type Info struct {
	Platform string
	Arch     string
}

func MapBinaryTarget(name string) Info {
	return Info{
		Platform: mapBinaryTargetToPlatform(name),
		Arch:     mapBinaryTargetToArch(name),
	}
}

func mapBinaryTargetToPlatform(name string) string {
	switch {
	case strings.Contains(name, "linux") ||
		strings.Contains(name, "debian") ||
		strings.Contains(name, "rhel") ||
		strings.Contains(name, "musl"):
		return "linux"
	case strings.Contains(name, "darwin"):
		return "darwin"
	case strings.Contains(name, "windows"):
		return "windows"
	default:
		return "linux"
	}
}

func mapBinaryTargetToArch(name string) string {
	switch {
	case strings.Contains(name, "arm64"):
		return "arm64"
	default:
		return "!arm64"
	}
}
