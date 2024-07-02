package instance

import (
	"fmt"
	"os"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func InstanceList(instance []*models.Instance, choice chan<- string) {
	itemsList := make([]list.Item, len(instance))

	items := map[string]string{}

	for i, r := range instance {
		items[r.Name] = r.Name
		itemsList[i] = selection.Item(r.Name)
	}

	m := selection.NewModel(itemsList, "Instance")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- items[p.Choice]
	}

}