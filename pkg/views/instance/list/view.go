package list

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: 6},
	{Title: "Name", Width: 12},
	{Title: "Provider", Width: 12},
	{Title: "Endpoint", Width: 26},
	{Title: "Status", Width: 12},
	{Title: "Auth Mode", Width: 12},
	{Title: "Description", Width: 12},
}

func ListInstance(instance []*models.Instance) {
	var rows []table.Row
	for _, regis := range instance {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", regis.ID),
			regis.Name,
			regis.Vendor,
			regis.Endpoint,
			regis.Status,
			regis.AuthMode,
			regis.Description,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
