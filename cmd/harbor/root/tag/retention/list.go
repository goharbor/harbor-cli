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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	policylist "github.com/goharbor/harbor-cli/pkg/views/retention/policy/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCommand() *cobra.Command {
	var retentionID int64

	cmd := &cobra.Command{
		Use:   "list [PROJECT_NAME]",
		Short: "display retention policy for a project",
		Long:  "retrieve and display retention policy configured for a project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := ""
			if len(args) > 0 {
				projectName = args[0]
			}

			if retentionID == 0 {
				if projectName == "" {
					var err error
					projectName, err = prompt.GetProjectNameFromUser()
					if err != nil {
						return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
					}
				}

				resolvedID, err := api.GetRetentionIDByProjectName(projectName)
				if err != nil {
					return fmt.Errorf("failed to resolve retention policy for project %q: %v", projectName, utils.ParseHarborErrorMsg(err))
				}
				retentionID = resolvedID
			}

			policy, err := api.GetRetentionPolicy(retentionID)
			if err != nil {
				return fmt.Errorf("failed to get retention policy %d: %v", retentionID, utils.ParseHarborErrorMsg(err))
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				if err := utils.PrintFormat(policy, formatFlag); err != nil {
					return err
				}
				return nil
			}

			if policy == nil {
				fmt.Println("No retention policy found.")
				return nil
			}

			policylist.ListPolicy(policy)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&retentionID, "retention-id", "", 0, "retention policy ID")

	return cmd
}
