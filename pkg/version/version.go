// Package version contains build-time version information for the application.
package version

import (
	"fmt"
)

// Version is the current application version.
var Version = ""

// BuildDate is the date when the application was built.
var BuildTime = ""

// CommitHash is the git commit hash of the build.
var CommitHash = ""

func PrintVersion() {
	fmt.Printf("Version: %s\nBuild time: %s\nCommitHash: %s\n", Version, BuildTime, CommitHash)
}
