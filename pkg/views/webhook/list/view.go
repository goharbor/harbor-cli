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
package list

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: 6},
	{Title: "Webhook Name", Width: 15},
	{Title: "Enabled", Width: 12},
	{Title: "Endpoint URL", Width: 40},
	{Title: "Notify Type", Width: 12},
	{Title: "Payload Format", Width: 15},
	{Title: "Creation Time", Width: 20},
}

func ListWebhooks(webhooks []*models.WebhookPolicy) {
	var rows []table.Row
	for _, webhook := range webhooks {
		var webhookEnabled string
		if webhook.Enabled {
			webhookEnabled = "True"
		} else {
			webhookEnabled = "False"
		}
		payloadFormat := "--"
		if len(webhook.Targets[0].PayloadFormat) != 0 {
			payloadFormat = string(webhook.Targets[0].PayloadFormat)
		}
		creationTime, _ := utils.FormatCreatedTime(webhook.CreationTime.String())
		rows = append(rows, table.Row{
			strconv.FormatInt(webhook.ID, 10),
			webhook.Name,
			webhookEnabled,
			webhook.Targets[0].Address,
			webhook.Targets[0].Type,
			payloadFormat,
			creationTime,
		})
	}
	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
