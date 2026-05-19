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
package policy

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func StartPolicyCommand() *cobra.Command {
	var isID bool

	cmd := &cobra.Command{
		Use:     "start [NAME|ID] [POLICY_NAME]",
		Short:   "Manually trigger a preheat policy",
		Long:    "Manually trigger a specific P2P preheat policy under a project",
		Example: `  harbor-cli project preheat policy start [NAME|ID] [POLICY_NAME]`,
		Args:    cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName, policyName string

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

			log.Debug("Manually triggering preheat policy...")
			location, err := api.StartPreheatPolicy(projectName, policyName)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "404" {
					return fmt.Errorf("preheat policy %s not found in project %s", policyName, projectName)
				}
				return fmt.Errorf("failed to manually trigger preheat policy: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Printf("Preheat policy '%s' manually triggered successfully from project '%s'\nExecution location: %s\n", policyName, projectName, location)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&isID, "id", false, "Manually trigger preheat policy by project id")

	return cmd
}
