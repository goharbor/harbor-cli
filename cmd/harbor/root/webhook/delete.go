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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func DeleteWebhookCmd() *cobra.Command {
	var (
		projectName string
		webhookId   int64 = -1
	)

	cmd := &cobra.Command{
		Use:   "delete [webhook-name]",
		Short: "Delete a webhook from a Harbor project",
		Long: `Delete a webhook from a specified Harbor project.
You can either specify the project name and webhook ID using flags,
pass the webhook name as an argument,
or interactively select a project and webhook if not provided.`,
		Example: `  # Delete by project and webhook ID
  harbor-cli webhook delete --project my-project --webhook 5

  # Delete by project and webhook name
  harbor-cli webhook delete my-webhook --project my-project

  # Fully interactive deletion
  harbor-cli webhook delete`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var webhookName string

			if len(args) > 0 {
				webhookName = args[0]
			}

			if webhookId != -1 && webhookName != "" {
				return fmt.Errorf("webhook ID and name cannot be provided together")
			}

			if (webhookId != -1 || webhookName != "") && projectName == "" {
				return fmt.Errorf("project name must be provided when specifying webhook ID or webhook name")
			}

			if projectName == "" {
				projectName = prompt.GetProjectNameFromUser()
			}

			if webhookName != "" {
				webhookId, err = api.GetWebhookID(projectName, webhookName)
				if err != nil {
					return fmt.Errorf("failed to get webhook ID for '%s': %w", webhookName, err)
				}
			}

			if webhookId == -1 {
				selectedWebhook, err := prompt.GetWebhookFromUser(projectName)
				if err != nil {
					return fmt.Errorf("failed to select webhook: %w", err)
				}
				webhookId = selectedWebhook.ID
			}

			if err := api.DeleteWebhook(projectName, webhookId); err != nil {
				return fmt.Errorf("failed to delete webhook: %w", err)
			}

			fmt.Printf("Webhook deleted successfully from project '%s'\n", projectName)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectName, "project", "", "", "Project name (required when providing webhook ID or name)")
	flags.Int64VarP(&webhookId, "webhook-id", "", -1, "Webhook ID (alternative to providing webhook name as argument)")

	return cmd
}
