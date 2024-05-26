package registry

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views"
)

func ListArtifacts(artifacts []*models.Artifact, choice chan<- string) {
	itemsList := make([]list.Item, len(artifacts))

	for i, a := range artifacts {
		itemsList[i] = views.Item(a.Digest)
	}

	m := views.NewModel(itemsList, "Artifact")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(views.Model); ok {
		choice <- p.Choice
	}

}
