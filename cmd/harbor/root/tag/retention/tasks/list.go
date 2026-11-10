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
	"github.com/goharbor/harbor-cli/pkg/utils"
	taskview "github.com/goharbor/harbor-cli/pkg/views/retention/tasks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCommand() *cobra.Command {
	var projectName string
	var projectID int
	var executionID int64

	cmd := &cobra.Command{
		Use:   "list [PROJECT_NAME]",
		Short: "List retention tasks for an execution",
		Long:  "List repository-level retention tasks for a specific retention execution",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				projectName = args[0]
			}

			if projectID != -1 && projectName != "" {
				return fmt.Errorf("cannot specify both --project-id and --project-name flags")
			}

			if projectID == -1 && projectName == "" {
				name, err := prompt.GetProjectNameFromUser()
				if err != nil {
					return err
				}
				projectName = name
			}

			projectIDStr := projectName
			isName := true
			if projectID != -1 {
				projectIDStr = strconv.Itoa(projectID)
				isName = false
			}

			retentionID, err := api.GetRetentionId(projectIDStr, isName)
			if err != nil {
				return fmt.Errorf("failed to resolve retention policy: %w", err)
			}

			if executionID == -1 {
				executionID = prompt.GetRetentionExecutionIDFromUser(retentionID)
			}
			if executionID == 0 {
				fmt.Println("No retention executions found")
				return nil
			}

			resp, err := api.ListRetentionTasks(retentionID, executionID)
			if err != nil {
				return fmt.Errorf("failed to list retention tasks: %w", err)
			}

			if viper.GetString("output-format") != "" {
				utils.PrintPayloadInJSONFormat(resp.Payload)
				return nil
			}

			taskview.ListRetentionTasks(resp.Payload)
			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().IntVarP(&projectID, "project-id", "i", -1, "Project ID")
	cmd.Flags().Int64VarP(&executionID, "execution-id", "e", -1, "Retention execution ID")

	return cmd
}
