package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the application version
	Version = "1.0.0"
	// BuildTime is the build timestamp
	BuildTime = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// GoVersion is the Go version used to build the application
	GoVersion = runtime.Version()
)

// Info returns version information
func Info() string {
	return fmt.Sprintf("Version: %s, BuildTime: %s, GitCommit: %s, GoVersion: %s",
		Version, BuildTime, GitCommit, GoVersion)
}

// Short returns a short version string
func Short() string {
	return fmt.Sprintf("v%s", Version)
}
