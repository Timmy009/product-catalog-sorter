package version

import (
	"fmt"
	"runtime"
)

// BuildInfo contains build-time information
type BuildInfo struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
	BuildTime  string `json:"build_time"`
	GoVersion  string `json:"go_version"`
}

// String returns a formatted string representation of build info
func (b BuildInfo) String() string {
	return fmt.Sprintf("Version: %s, Commit: %s, Built: %s, Go: %s",
		b.Version, b.CommitHash, b.BuildTime, b.GoVersion)
}

// GetRuntimeInfo returns current runtime information
func GetRuntimeInfo() RuntimeInfo {
	return RuntimeInfo{
		GoVersion:    runtime.Version(),
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
	}
}

// RuntimeInfo contains runtime information
type RuntimeInfo struct {
	GoVersion    string `json:"go_version"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
}

// String returns a formatted string representation of runtime info
func (r RuntimeInfo) String() string {
	return fmt.Sprintf("Go: %s, OS: %s, Arch: %s, CPUs: %d, Goroutines: %d",
		r.GoVersion, r.GOOS, r.GOARCH, r.NumCPU, r.NumGoroutine)
}
