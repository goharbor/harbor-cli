package retention

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/retention/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListRetentionRulesCommand() *cobra.Command {
	var projectName string
	var projectID int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tag retention rules of a project",
		Long: `Retrieve and display the tag retention rules for a specific project in Harbor.

Tag retention rules define policies for automatically keeping or deleting image tags 
within a project. Using this command, you can view the currently configured 
retention rules.

Usage:
  - Specify the project **either by name or by ID**, but not both.
  - If neither is provided, you will be prompted to select a project.
  - The rules will be displayed in a formatted output.

Examples:
  # List retention rules using project name
  harbor tag retention list --project-name my-project

  # List retention rules using project ID
  harbor tag retention list --project-id 42

  # List retention rules in JSON format
  harbor tag retention list --project-name my-project --output-format json`,
		Args: cobra.NoArgs,
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
			resp, err := api.ListRetention(retentionID)
			if err != nil {
				return fmt.Errorf("failed to list retention rules: %w", err)
			}
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				utils.PrintPayloadInJSONFormat(resp)
				return nil
			}

			list.ListRetentionRules(resp.Payload.Rules)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().IntVarP(&projectID, "project-id", "i", -1, "Project ID")

	return cmd
}
