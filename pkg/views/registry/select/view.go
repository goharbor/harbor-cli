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

func RegistryListTypes(registries []string, choice chan<- string) {

	itemsList := make([]list.Item, len(registries))

	items := map[string]string{}

	for i, reg := range registries {
		items[reg] = reg
		itemsList[i] = item(reg)
	}

	const defaultWidth = 20

	l := list.New(itemsList, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select a Registry"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = views.TitleStyle
	l.Styles.PaginationStyle = views.PaginationStyle
	l.Styles.HelpStyle = views.HelpStyle

	m := model{list: l}

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(model); ok {
		choice <- items[p.choice]
	}

}

