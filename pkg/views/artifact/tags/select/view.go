package registry

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func ListTags(tag []*models.Tag, choice chan<- string) {
	itemsList := make([]list.Item, len(tag))

	for i, t := range tag {
		itemsList[i] = selection.Item(t.Name)
	}

	m := selection.NewModel(itemsList, "Tag")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- p.Choice
	}
}
