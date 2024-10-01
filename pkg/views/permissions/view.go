package permissions

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var systemColumns = []table.Column{
	{Title: "Resource", Width: 30},
	{Title: "Action", Width: 20},
}

var projectColumns = []table.Column{
	{Title: "Resource", Width: 30},
	{Title: "Action", Width: 20},
}

func PrintPermissions(perms *models.Permissions) {
	var systemRows []table.Row
	for _, p := range perms.System {
		systemRows = append(systemRows, table.Row{
			p.Resource,
			p.Action,
		})
	}

	systemModel := tablelist.NewModel(systemColumns, systemRows, len(systemRows))
	var projectRows []table.Row
	for _, p := range perms.Project {
		projectRows = append(projectRows, table.Row{
			p.Resource,
			p.Action,
		})
	}

	projectModel := tablelist.NewModel(projectColumns, projectRows, len(projectRows))

	fmt.Println("System Permissions:")
	if _, err := tea.NewProgram(systemModel).Run(); err != nil {
		fmt.Println("Error running system permissions table:", err)
		os.Exit(1)
	}

	fmt.Println("\nProject Permissions:")
	if _, err := tea.NewProgram(projectModel).Run(); err != nil {
		fmt.Println("Error running project permissions table:", err)
		os.Exit(1)
	}
}
