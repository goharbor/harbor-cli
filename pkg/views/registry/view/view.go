package view

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: 5},
	{Title: "Name", Width: 10},
	{Title: "Status", Width: 10},
	{Title: "Endpoint URL", Width: 20},
	{Title: "Provider", Width: 15},
	{Title: "Creation Time", Width: 15},
	{Title: "Description", Width: 20},
}

func ViewRegistry(registry *models.Registry) {
	var rows []table.Row
	createdTime, _ := utils.FormatCreatedTime(registry.CreationTime.String())
	rows = append(rows, table.Row{
		fmt.Sprintf("%d", registry.ID),
		registry.Name,
		registry.Status,
		registry.URL,
		registry.Type,
		createdTime,
		registry.Description,
	})

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
