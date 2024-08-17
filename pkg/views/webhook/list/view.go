package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"os"
	"strconv"
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
