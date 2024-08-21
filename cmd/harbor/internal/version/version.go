// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
