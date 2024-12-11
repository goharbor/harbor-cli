package list

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
	{Title: "ID", Width: 6},
	{Title: "Name", Width: 12},
	{Title: "Color", Width: 12},
	{Title: "Description", Width: 18},
	{Title: "Creation Time", Width: 24},
}

func ListLabels(labels []*models.Label) {
	var rows []table.Row
	for _, regis := range labels {
		createdTime, _ := utils.FormatCreatedTime(regis.CreationTime.String())
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", regis.ID),
			regis.Name,
			regis.Color,
			regis.Description,
			createdTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
