package robot

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/multiselect"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func ListPermissions(perms *models.Permissions, ch chan<- []models.Permission) {
	permissions := perms.Project
	choices := []models.Permission{}

	// Iterate over permissions and append each item to choices
	for _, perm := range permissions {
		choices = append(choices, *perm)
	}

	selects := &[]models.Permission{}

	m := multiselect.NewModel(choices, selects)

	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
	}
	// Get selected permissions
	ch <- *selects
}

func ListRobot(robots []*models.Robot, choice chan<- int64) {
	itemsList := make([]list.Item, len(robots))

	items := map[string]int64{}

	for i, r := range robots {
		items[r.Name] = r.ID
		itemsList[i] = selection.Item(r.Name)
	}

	m := selection.NewModel(itemsList, "Robot")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- items[p.Choice]
	}
}
