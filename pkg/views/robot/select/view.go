package robot

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/multiselect"
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

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
	}
	// Get selected permissions
	ch <- *selects
}
