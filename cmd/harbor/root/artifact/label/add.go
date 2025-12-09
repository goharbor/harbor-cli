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

// AddLabelArtifactCommmand adds a label to an artifact
func AddLabelArtifactCommmand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Attach a label to an artifact in a Harbor project repository",
		Long: `Attach an existing label to a specific artifact identified by <project>/<repository>:<reference>.
You can specify the artifact and label directly as arguments, or interactively select them if arguments are omitted.

Examples:
  # Add a label to an artifact using project/repo:reference and label name
  harbor artifact label add myproject/myrepo@sha256:abcdef1234567890 dev

  # Prompt-based label selection for an artifact
  harbor artifact label add library/nginx:1.21

  # Fully interactive mode (prompt for everything)
  harbor artifact label add
`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				projectName, repoName, reference string
				labelName                        string
				labelID                          int64
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
				labelName = args[1]
				labelID, err = api.GetLabelIdByName(labelName, api.ListFlags{})
				if err != nil {
					return fmt.Errorf("failed to get label id: %v", utils.ParseHarborErrorMsg(err))
				}
			} else {
				labelID, err = prompt.GetLabelIdFromUser(api.ListFlags{})
				if err != nil {
					return fmt.Errorf("failed to get label id: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			label := api.GetLabel(labelID)

			if _, err := api.AddLabelArtifact(projectName, repoName, reference, label); err != nil {
				return fmt.Errorf("failed to add label to artifact: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Printf("Label '%s' added to artifact %s/%s@%s\n", label.Name, projectName, repoName, reference)
			return nil
		},
	}

	return cmd
}
