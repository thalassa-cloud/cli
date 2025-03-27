package version

import (
	"fmt"
)

type VersionInfo struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

var versionInfo *VersionInfo

// Init initializes the version information
func Init(version, commit, date, builtBy string) {
	versionInfo = &VersionInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	}
}

// PrintVersion
func PrintVersion() {
	if versionInfo == nil {
		fmt.Println("Version information is not initialized")
		return
	}
	fmt.Println(versionInfo)
}

// Commit returns the git commit
func Commit() string {
	if versionInfo == nil {
		return "unknown"
	}
	return versionInfo.Commit
}

// Version returns the application version
func Version() string {
	if versionInfo == nil {
		return "unknown"
	}
	return versionInfo.Version
}

// BuildDate returns the build date
func BuildDate() string {
	if versionInfo == nil {
		return "unknown"
	}
	return versionInfo.Date
}

// BuiltBy returns who built the application
func BuiltBy() string {
	if versionInfo == nil {
		return "unknown"
	}
	return versionInfo.BuiltBy
}

// String returns the formatted version information
func (v *VersionInfo) String() string {
	return fmt.Sprintf("Version information\nVersion: %s, Commit: %s\nBuild date: %s, Built by: %s\n", v.Version, v.Commit, v.Date, v.BuiltBy)
}
