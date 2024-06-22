package version

import (
	"fmt"
	"runtime/debug"
	"time"
)

const (
	VersionMajor = 0          // Major version for stable releases
	VersionMinor = 1          // Minor version for stable releases
	VersionPatch = 0          // Patch version for stable releases
	VersionMeta  = "unstable" // Version metadata to append to the version string (e.g. "unstable" or "stable" or "beta")
)

// VersionInfo holds the commit hash and commit hash date.
type VersionInfo struct {
	CommitHash              string
	CommitHashDate          string
	FormattedCommitHashDate string
}

// GetVersionInfo retrieves the current commit hash and commit hash date of the application.
func GetVersionInfo() VersionInfo {
	var versionInfo VersionInfo
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				versionInfo.CommitHash = setting.Value
			case "vcs.time":
				versionInfo.CommitHashDate = setting.Value
				versionInfo.FormattedCommitHashDate = formatDate(setting.Value)
			}
		}
	}

	return versionInfo
}

// Version is the current version of the application.
var Version = func() string {
	vi := GetVersionInfo()
	return fmt.Sprintf("%d.%d.%d-%s (%s %s)", VersionMajor, VersionMinor, VersionPatch, VersionMeta, vi.CommitHash, vi.FormattedCommitHashDate)
}()

// formatDate formats a date string from RFC3339 to yyyy-mm-dd.
func formatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}
