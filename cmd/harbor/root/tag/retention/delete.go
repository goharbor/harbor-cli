package retention

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteRetentionPolicyCommand() *cobra.Command {
	var projectName string
	var projectID int

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a tag retention policy for a project",
		Long: `Delete an existing tag retention policy from a project in Harbor.

Usage:
  - You can specify the project **either by name or by ID**, but not both.
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
				return fmt.Errorf("%w", err)
			}
			if err := api.DeleteRetention(retentionID); err != nil {
				return fmt.Errorf("failed to delete retention rule: %w", err)
			}

			log.Info("Retention Policy deleted successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().IntVarP(&projectID, "project-id", "i", -1, "Project ID")

	return cmd
}
