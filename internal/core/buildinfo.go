package core

// BuildInfo holds version information about the build
// (moved from internal/build/buildinfo.go)
type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

// DefaultBuildInfo provides default build information when not set during compilation
var DefaultBuildInfo = BuildInfo{
	Version: "dev",
	Commit:  "none",
	Date:    "unknown",
}
