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
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteWebhookCmd() *cobra.Command {
	var projectName string
	var webhookId string
	var webhookIdInt int64
	var selectedWebhook models.WebhookPolicy

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a webhook from a Harbor project",
		Long: `This command deletes a webhook from the specified Harbor project.
You can either specify the project name and webhook ID directly using flags,
or interactively select a project and webhook if not provided.`,
		Example: `  # Delete a webhook by specifying the project and webhook ID
  harbor-cli webhook delete --project my-project --webhook 5

  # Delete a webhook by selecting the project and webhook interactively
  harbor-cli webhook delete`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if projectName != "" && webhookId != "" {
				webhookIdInt, err = strconv.ParseInt(webhookId, 10, 64)
				if err != nil {
					log.Errorf("failed to parse webhook id: %v", err)
					return
				}
			} else {
				projectName = prompt.GetProjectNameFromUser()
				selectedWebhook = prompt.GetWebhookFromUser(projectName)
				webhookIdInt = selectedWebhook.ID
			}
			err = api.DeleteWebhook(projectName, webhookIdInt)
			if err != nil {
				log.Errorf("failed to delete webhook: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectName, "project", "", "", "Project Name")
	flags.StringVarP(&webhookId, "webhook", "", "", "Webhook ID")

	return cmd
}
