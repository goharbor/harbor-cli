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
package project

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/summary"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SummaryCommand creates a new harbor project summary command
func SummaryCommand() *cobra.Command {
	var isID bool
	cmd := &cobra.Command{
		Use:     "summary [NAME|ID]",
		Short:   "Get summary of a project",
		Long:    "Get summary of a project by name or ID. If no arguments are provided, it will prompt for the project name. Use --id to specify the project ID instead of the name.",
		Example: "harbor project summary my-project or harbor project summary 1 --id",
		Args:    cobra.MaximumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			var err error

			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			projectData, err := api.GetProject(projectName, isID)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "404" {
					return fmt.Errorf("project %s does not exist", projectName)
				}
				return fmt.Errorf("failed to get project details: %v", utils.ParseHarborErrorMsg(err))
			}

			projectSummary, err := api.GetProjectSummary(projectName, isID)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "404" {
					return fmt.Errorf("project %s does not exist", projectName)
				}
				return fmt.Errorf("failed to get project summary: %v", utils.ParseHarborErrorMsg(err))
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				combined := map[string]interface{}{
					"project": projectData.Payload,
					"summary": projectSummary.Payload,
				}
				err = utils.PrintFormat(combined, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				err = summary.ViewProjectSummary(projectData.Payload, projectSummary.Payload)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&isID, "id", false, "Get project by id")

	return cmd
}
