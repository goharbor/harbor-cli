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
package label

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/label/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DelLabelArtifactCommmand delete label command to artifact
func ListLabelArtifactCommmand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Display labels attached to a specific artifact",
		Long: `This command lists all labels currently associated with a specific artifact in a Harbor project repository.
You can provide the artifact reference in the format <project>/<repository>:<reference> (where reference is either a tag or a digest).
If the reference is not provided as an argument, the command will prompt you to select the project, repository, and artifact.

Supports output formatting such as JSON or YAML using the --output (-o) flag.`,
		Example: `  # List labels for a tagged artifact
  harbor artifact label list library/nginx:latest

  # List labels for an artifact by digest
  harbor artifact label list myproject/myrepo@sha256:abc123...

  # Prompt-based interactive selection of artifact
  harbor artifact label list

  # Output in JSON format
  harbor artifact label list library/nginx:1.21 -o json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName, repoName, reference string
			var artifact *artifact.GetArtifactOK
			getLabel := true
			if len(args) > 0 {
				projectName, repoName, reference, err = utils.ParseProjectRepoReference(args[0])
				if err != nil {
					return fmt.Errorf("failed to parse project/repo/reference: %v", err)
				}
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
				repoName = prompt.GetRepoNameFromUser(projectName)
				reference = prompt.GetReferenceFromUser(repoName, projectName)
			}

			if reference == "" {
				if len(args) > 0 {
					return fmt.Errorf("Invalid artifact reference format: %s", args[0])
				} else {
					return fmt.Errorf("Invalid artifact reference format: no arguments provided")
				}
			}

			artifact, err = api.ViewArtifact(projectName, repoName, reference, getLabel)

			if err != nil || artifact == nil {
				return fmt.Errorf("failed to get info of an artifact: %v", utils.ParseHarborErrorMsg(err))
			}
			labelList := artifact.Payload.Labels
			if len(labelList) == 0 {
				fmt.Printf("No labels found for artifact %s/%s@%s", projectName, repoName, reference)
				return nil
			}
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(labelList, formatFlag)
				if err != nil {
					return err
				}
			} else {
				list.ListLabels(labelList)
			}
			return nil
		},
	}

	return cmd
}
