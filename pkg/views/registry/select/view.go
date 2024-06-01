package registry

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func RegistryList(registry []*models.Registry, choice chan<- int64) {
	itemsList := make([]list.Item, len(registry))

	items := map[string]int64{}

	for i, r := range registry {
		items[r.Name] = r.ID
		itemsList[i] = selection.Item(r.Name)
	}

	m := selection.NewModel(itemsList, "Registry")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- items[p.Choice]
	}

}
