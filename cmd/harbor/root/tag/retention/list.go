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
		Use:     "list",
		Short:   "List retention execution of the project",
		Args:    cobra.NoArgs,
		Example: `harbor retention list --project-name myproject`,
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
				return fmt.Errorf("No retention policy exists for this project")
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
	cmd.Flags().IntVarP(&projectID, "project-id", "i", 0, "Project ID")

	return cmd
}
