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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
)

func DeleteWebhookCmd() *cobra.Command {

	var projectName string
	var webhookId string
	var webhookIdInt int64
	var selectedWebhook models.WebhookPolicy

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "webhook delete",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if projectName != "" && webhookId != "" {
				webhookIdInt, err = strconv.ParseInt(webhookId, 10, 64)
				err = api.DeleteWebhook(projectName, webhookIdInt)
			} else {
				projectName = prompt.GetProjectNameFromUser()
				selectedWebhook = prompt.GetWebhookFromUser(projectName)
				err = api.DeleteWebhook(projectName, selectedWebhook.ID)

			}
			if err != nil {
				log.Errorf("failed to delete webhook: %v", err)
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectName, "project", "", "", "Project Name")
	flags.StringVarP(&webhookId, "webhook", "", "", "Webhook Id")

	return cmd
}
