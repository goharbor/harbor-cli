package version

import "runtime/debug"

var (
	Version        = "0.1.0"
	GitCommit      = ""
	BuildTime      = ""
	ReleaseChannel = "dev"
	GoVersion      = ""
	OS             = func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "GOOS" {
					return setting.Value
				}
			}
		}

		return ""
	}
	Arch = func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "GOARCH" {
					return setting.Value
				}
			}
		}

		return ""
	}
	System = OS() + "/" + Arch()
)
