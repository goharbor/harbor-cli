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
		Short: "Delete retention policy for a project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectID != 0 && projectName != "" {
				return fmt.Errorf("Cannot specify both --project-id and --project-name flags")
			}

			if projectID == 0 && projectName == "" {
				projectName = prompt.GetProjectNameFromUser()
			}

			projectIDStr := ""
			isName := true
			if projectID != 0 {
				projectIDStr = strconv.Itoa(projectID)
				isName = false
			} else {
				projectIDStr = projectName
			}

			retentionID, err := api.GetRetentionId(projectIDStr, isName)
			if err != nil {
				return fmt.Errorf("No retention policy exists for this project: %w", err)
			}
			if err := api.DeleteRetention((retentionID)); err != nil {
				return fmt.Errorf("failed to delete retention rule: %w", err)
			}

			log.Info("Retention Policy deleted successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().IntVarP(&projectID, "project-id", "i", 0, "Project ID")

	return cmd
}
