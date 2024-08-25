package project

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func ProjectList(project []*models.Project, choice chan<- string) {
	items := make([]list.Item, len(project))
	for i, p := range project {
		items[i] = selection.Item(p.Name)
	}

	m := selection.NewModel(items, "Project")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- p.Choice
	}

}

func ProjectIdList(project []*models.Project, choice chan<- int32) {
	itemsList := make([]list.Item, len(project))
	items := map[string]int32{}
	for i, p := range project {
		items[p.Name] = p.ProjectID
		itemsList[i] = selection.Item(p.Name)
	}

	m := selection.NewModel(itemsList, "Project")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- items[p.Choice]
	}

}
