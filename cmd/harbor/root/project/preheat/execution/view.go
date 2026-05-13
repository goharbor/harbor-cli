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
package execution

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/preheat/execution/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ViewExecutionCommand() *cobra.Command {
	var isID bool

	cmd := &cobra.Command{
		Use:     "view [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID]",
		Short:   "View preheat execution details",
		Long:    "Get details of a specific P2P preheat execution under a project",
		Example: `  harbor-cli project preheat execution view [NAME|ID] [POLICY_NAME] [EXECUTION_ID]`,
		Args:    cobra.MaximumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName, policyName string
			var executionID int64

			if isID && len(args) == 0 {
				return fmt.Errorf("project ID must be provided when using --id")
			}

			if len(args) >= 1 {
				log.Debugf("Project name provided: %s", args[0])
				projectName = args[0]
			} else {
				log.Debug("No project name provided, prompting user")
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
			}
			if isID {
				project, err := api.GetProject(projectName, true)
				if err != nil {
					return fmt.Errorf("failed to resolve project ID: %v", utils.ParseHarborErrorMsg(err))
				}
				projectName = project.Payload.Name
			}

			if len(args) >= 2 {
				log.Debugf("Policy name provided: %s", args[1])
				policyName = args[1]
			} else {
				log.Debug("No policy name provided, prompting user")
				policyName, err = prompt.GetPreheatPolicyNameFromUser(projectName)
				if err != nil {
					return fmt.Errorf("failed to get policy name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			if len(args) >= 3 {
				log.Debugf("Execution ID provided: %s", args[2])
				executionID, err = strconv.ParseInt(args[2], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid execution ID %q: %v", args[2], err)
				}
			} else {
				log.Debug("No execution ID provided, prompting user")
				executionID, err = prompt.GetPreheatPolicyExecIDFromUser(projectName, policyName)
				if err != nil {
					return fmt.Errorf("failed to get execution id: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			log.Debug("Fetching preheat execution details...")
			resp, err := api.GetPreheatExecution(projectName, policyName, executionID)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "404" {
					return fmt.Errorf("no execution found for execution ID %d in policy %s in project %s", executionID, policyName, projectName)
				}
				return fmt.Errorf("failed to view preheat execution: %v", utils.ParseHarborErrorMsg(err))
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(resp.Payload, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				view.ViewExecution(resp.Payload)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&isID, "id", false, "Get preheat policy execution by project id")

	return cmd
}
