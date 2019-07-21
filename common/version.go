package common

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// This is the current version of the client
	CurrentVersion = Version{
		Major:  0,
		Minor:  1,
		Patch:  1,
		Suffix: "",
	}

	// commitHash contains the current Git revision. Use Go Releaser to make sure this gets set.
	commitHash string

	// buildDate contains the date of the current build.
	buildDate string
)

type Version struct {
	// Major version
	Major uint32

	// Minor version
	Minor uint32

	// Patch version
	Patch uint32

	// Suffix used in version string
	// Will be blank for release versions
	Suffix string
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d%s", v.Major, v.Minor, v.Patch, v.Suffix)
}

func BuildVersionString() string {
	program := "Portainer Stack Utils"

	version := "v" + CurrentVersion.String()

	if commitHash != "" {
		version += "-" + strings.ToUpper(commitHash)
	}

	osArch := runtime.GOOS + "/" + runtime.GOARCH

	date := buildDate
	if date == "" {
		date = "unknown"
	}

	return fmt.Sprintf("%s %s %s BuildDate: %s", program, version, osArch, date)
}
