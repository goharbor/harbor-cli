package project

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
	"os"
)

func WebhookList(webhooks []*models.WebhookPolicy, choice chan<- models.WebhookPolicy) {
	itemsList := make([]list.Item, len(webhooks))

	for i, item := range webhooks {
		itemsList[i] = selection.Item(item.Name)
	}

	m := selection.NewModel(itemsList, "Webhook")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		for _, webhook := range webhooks {
			if webhook.Name == p.Choice {
				choice <- *webhook
				break
			}
		}
	}

}
