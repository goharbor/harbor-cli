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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// DelLabelArtifactCommmand deletes a label from an artifact
func DelLabelArtifactCommmand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "Detach a label from an artifact in a Harbor project repository",
		Long: `Remove an existing label from a specific artifact identified by <project>/<repository>:<reference>.
You can provide the artifact and label name as arguments, or choose them interactively if not specified.

Examples:
  # Remove a label by specifying artifact and label name
  harbor artifact label delete library/nginx:latest stable

  # Prompt-based label selection for a specific artifact
  harbor artifact label del library/nginx:1.21

  # Fully interactive mode (prompt for project, repo, reference, and label)
  harbor artifact label delete

  # Remove a label from an artifact identified by digest
  harbor artifact label del myproject/myrepo@sha256:abcdef1234567890 qa-label`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				projectName, repoName, reference string
				labelID                          int64 = -1
				err                              error
			)

			if len(args) >= 1 {
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

			if len(args) == 2 {
				labelName := args[1]
				labelID, err = api.GetLabelIdByName(labelName)
				if err != nil {
					return fmt.Errorf("failed to get label id: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			if labelID == -1 {
				artifact, err := api.ViewArtifact(projectName, repoName, reference, true)
				if err != nil || artifact == nil {
					return fmt.Errorf("failed to get artifact info: %v", utils.ParseHarborErrorMsg(err))
				}

				labels := artifact.Payload.Labels
				if len(labels) == 0 {
					fmt.Printf("No labels found for artifact %s/%s@%s\n", projectName, repoName, reference)
					return nil
				}
				labelID = prompt.GetLabelIdFromUser(labels)
			}

			if _, err := api.RemoveLabelArtifact(projectName, repoName, reference, labelID); err != nil {
				return fmt.Errorf("failed to remove label from artifact: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Println("Label removed successfully")
			return nil
		},
	}

	return cmd
}
