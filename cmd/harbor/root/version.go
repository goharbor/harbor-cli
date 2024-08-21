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

package root

import (
	"fmt"

	"github.com/goharbor/harbor-cli/cmd/harbor/internal/version"
	"github.com/spf13/cobra"
)

// versionCommand represents the version command
func versionCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Version of Harbor CLI",
		Long:    `Get Harbor CLI version, git commit, go version, build time, release channel, os/arch, etc.`,
		Example: `  harbor version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion()
		},
	}

	return cmd
}

func runVersion() error {
	fmt.Printf("Version:      %s\n", version.Version)
	fmt.Printf("Go version:   %s\n", version.GoVersion)
	fmt.Printf("Git commit:   %s\n", version.GitCommit)
	fmt.Printf("Built:        %s\n", version.BuildTime)
	fmt.Printf("OS/Arch:      %s\n", version.System)

	return nil
}
