// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package artifact

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/artifact/view"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ViewArtifactCommmand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "Get information of an artifact",
		Long:    `Get information of an artifact`,
		Example: `harbor artifact view <project>/<repository>:<tag> OR harbor artifact view <project>/<repository>@<digest>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName, repoName, reference string
			var artifact *artifact.GetArtifactOK

			if len(args) > 0 {
				projectName, repoName, reference, err = utils.ParseProjectRepoReference(args[0])
				if err != nil {
					return fmt.Errorf("failed to parse project/repo/reference: %w", err)
				}
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %w", err)
				}
				repoName = prompt.GetRepoNameFromUser(projectName)
				reference = prompt.GetReferenceFromUser(repoName, projectName)
			}

			if reference == "" {
				if len(args) > 0 {
					return fmt.Errorf("invalid artifact reference format: %s", args[0])
				} else {
					return fmt.Errorf("invalid artifact reference format: no arguments provided")
				}
			}

			artifact, err = api.ViewArtifact(projectName, repoName, reference, false)

			if err != nil {
				return fmt.Errorf("failed to get info of an artifact: %w", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(artifact, FormatFlag)
				if err != nil {
					return fmt.Errorf("failed to print format: %w", err)
				}
			} else {
				view.ViewArtifact(artifact.Payload)
			}
			return nil
		},
	}

	return cmd
}
