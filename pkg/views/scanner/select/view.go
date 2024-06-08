package scanner

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
	"os"
)

func ScannerList(scanners []*models.ScannerRegistration, choice chan<- string) {
	itemsList := make([]list.Item, len(scanners))

	items := map[string]string{}

	for i, s := range scanners {
		items[s.Name] = s.UUID
		itemsList[i] = selection.Item(s.Name)
	}

	m := selection.NewModel(itemsList, "Scanner")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- items[p.Choice]
	}
}
