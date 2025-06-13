// Package version contains build-time version information for the application.
package version

import (
	"fmt"
)

// version is the current application version.
var version = ""

// buildDate is the date when the application was built.
var buildTime = ""

// commitHash is the git commit hash of the build.
var commitHash = ""

func PrintVersion() {
	fmt.Printf("Version: %s\nBuild time: %s\nCommitHash: %s\n", version, buildTime, commitHash)
}
