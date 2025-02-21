package envtpl

import (
	"fmt"
	"runtime"
)

const (
	AppName         = "envtpl"
	AppVersionMajor = 2
	AppVersionMinor = 0
	AppVersionPatch = 0
	AppVersionBuild = ""
)

// version build metadata set by build flag:
//     go build -ldflags "-X main.AppVersionMetadata $(date -u +%s)"
var AppVersionMetadata string

func Version() string {
	// major.minor.patch[-prerelease+buildmetadata]
	// optional version suffix format is "-(pre-release-version)+(build-metadata)"
	suffix := ""

	if AppVersionBuild != "" {
		suffix += "-" + AppVersionBuild
	}

	if AppVersionMetadata != "" {
		suffix += "-" + AppVersionMetadata
	}

	return fmt.Sprintf("%s %d.%d.%d%s (Go runtime %s).\nCopyright (c) 2016-2021, Tony Pujals.",
		AppName, AppVersionMajor, AppVersionMinor, AppVersionPatch, suffix, runtime.Version())
}
