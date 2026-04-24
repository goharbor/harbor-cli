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
	"github.com/spf13/cobra"
)

func RunCommand() *cobra.Command {
	var retentionID int64
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "run [PROJECT_NAME]",
		Short: "run retention policy",
		Long:  "trigger retention execution for a project policy",
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

			location, err := api.TriggerRetentionExecution(retentionID, dryRun)
			if err != nil {
				return fmt.Errorf("failed to trigger retention execution: %v", utils.ParseHarborErrorMsg(err))
			}

			if location != "" {
				fmt.Printf("Retention execution triggered successfully: %s\n", location)
				return nil
			}

			fmt.Println("Retention execution triggered successfully")
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&retentionID, "retention-id", "", 0, "retention policy ID")
	flags.BoolVarP(&dryRun, "dry-run", "", false, "trigger dry-run execution without deleting artifacts")

	return cmd
}
