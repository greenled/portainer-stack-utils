package common

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// This is the current version of the client. It is set by goreleaser.
	version string

	// The program name
	programName = "Portainer Stack Utils"

	// commitHash contains the current Git revision. Use Go Releaser to make sure this gets set.
	commitHash string

	// buildDate contains the date of the current build.
	buildDate string
)

func BuildVersionString() string {
	osArch := runtime.GOOS + "/" + runtime.GOARCH

	if version == "" {
		return fmt.Sprintf("%s SNAPSHOT %s", programName, osArch)
	}

	if commitHash != "" {
		version += "+" + strings.ToUpper(commitHash)
	}

	return fmt.Sprintf("%s %s %s BuildDate: %s", programName, version, osArch, buildDate)
}
