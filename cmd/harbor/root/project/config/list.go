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
package config

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/config/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListProjectConfigCmd() *cobra.Command {
	var err error
	var projectNameorID string
	cmd := &cobra.Command{
		Use:   "list  [project_name]",
		Short: "List configuration of a Harbor project by name or ID",
		Long: `Display the configuration metadata of a Harbor project specified by its name or ID.

If no project name or ID is provided as an argument, you will be prompted to select a project interactively.

You can use the global flag '--output-format' to specify the output format, e.g. 'json' or 'yaml', for machine-readable output.

Examples:

  # List configuration of project 'myproject' by name
  harbor-cli project config list myproject

  # List configuration of project with ID '123'
  harbor-cli project config list 123

  # Run interactively (prompt to select project)
  harbor-cli project config list

  # List config in JSON format
  harbor-cli project config list myproject --output-format json
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				projectNameorID, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
				isID = false
			} else {
				projectNameorID = args[0]
			}
			response, err := api.ListConfig(isID, projectNameorID)
			if err != nil {
				return fmt.Errorf("failed to list metadata: %v", utils.ParseHarborErrorMsg(err))
			}
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(response.Payload, formatFlag)
				if err != nil {
					return err
				}
			} else {
				list.ListConfig(response.Payload)
			}
			return nil
		},
	}
	return cmd
}
