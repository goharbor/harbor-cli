package project

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views"
)

func RepositoryList(repos []*models.Repository, choice chan<- string) {
	itemsList := make([]list.Item, len(repos))

	for i, r := range repos {
		split := strings.Split(r.Name, "/")
		itemsList[i] = views.Item(strings.Join(split[1:], "/"))
	}

	m := views.NewModel(itemsList, "Repository")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(views.Model); ok {
		choice <- p.Choice
	}

}
