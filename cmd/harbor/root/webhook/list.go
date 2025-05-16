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
package webhook

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/webhook"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	webhookViews "github.com/goharbor/harbor-cli/pkg/views/webhook/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListWebhookCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [PROJECT_NAME]",
		Short: "List all webhook policies for a Harbor project",
		Long: `This command retrieves and displays all webhook policies associated with a Harbor project.

You can either specify the project name directly as an argument or use the interactive prompt to select a project.
Use the '--output-format' flag for raw JSON output.`,
		Example: `  # List webhooks for a specific project
  harbor-cli webhook list my-project

  # List webhooks interactively by selecting the project
  harbor-cli webhook list

  # Output in JSON format
  harbor-cli webhook list my-project --output-format=json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var resp webhook.ListWebhookPoliciesOfProjectOK
			var projectName string

			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return err
				}
			}

			resp, err = api.ListWebhooks(projectName)
			if err != nil {
				return fmt.Errorf("failed to list webhooks: %v", err)
			}
			if len(resp.Payload) == 0 {
				fmt.Printf("No webhooks found for project %s\n", projectName)
				return nil
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(resp, FormatFlag)
				if err != nil {
					return fmt.Errorf("failed to print in %s format: %v", FormatFlag, err)
				}
			}
			webhookViews.ListWebhooks(resp.Payload)
			return nil
		},
	}
	return cmd
}
