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
package retention

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteRetentionRuleCommand() *cobra.Command {
	var projectName string
	var projectID int

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a tag retention policy for a project",
		Long: `Delete an existing tag retention policy from a project.

Usage:
  - You can specify the project either by name or by ID, but not both.
  - If neither is provided, you will be prompted to select a project.
  - The command retrieves the retention policy ID and deletes it.

Examples:
  # Delete retention policy using project name
  harbor tag retention delete --project-name my-project

  # Delete retention policy using project ID
  harbor tag retention delete --project-id 42`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID != -1 && projectName != "" {
				return fmt.Errorf("Cannot specify both --project-id and --project-name flags")
			}

			if projectID == -1 && projectName == "" {
				projectName = prompt.GetProjectNameFromUser()
			}

			projectIDStr := ""
			isName := true
			if projectID != -1 {
				projectIDStr = strconv.Itoa(projectID)
				isName = false
			} else {
				projectIDStr = projectName
			}

			retentionID, err := api.GetRetentionId(projectIDStr, isName)
			if err != nil {
				if err.Error() == "No retention policy exists for this project" {
					log.Info("No retention policy exists for this project")
					return nil
				} else {
					return fmt.Errorf("Error retrieving retention policy ID: %w", err)
				}

			}
			retentionIndex := prompt.GetRetentionTagRule(retentionID)
			err = api.DeleteRetention(projectName, int(retentionIndex))

			if err != nil {
				return fmt.Errorf("%w", err)
			}
			log.Info("Retention Rule deleted successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().IntVarP(&projectID, "project-id", "i", -1, "Project ID")

	return cmd
}
