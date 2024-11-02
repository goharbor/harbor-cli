package delete

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func LabelList(label []*models.Label, choice chan<- int64) {
	itemsList := make([]list.Item, len(label))

	items := map[string]int64{}

	for i, m := range label {
		items[m.Name] = m.ID
		itemsList[i] = selection.Item(m.Name)
	}

	m := selection.NewModel(itemsList, "Label")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- items[p.Choice]
	}

}
