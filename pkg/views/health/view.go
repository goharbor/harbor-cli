package health

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/health"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Component", Width: 20},
	{Title: "Status", Width: 30},
}

type HealthModel struct {
	tableModel tablelist.Model
}

func styleStatus(status string) string {
	if status == "healthy" {
		return views.GreenStyle.Render(status)
	}
	return views.RedStyle.Render(status)
}

func NewHealthModel(status *health.GetHealthOK) HealthModel {
	var rows []table.Row
	for _, component := range status.Payload.Components {
		rows = append(rows, table.Row{
			component.Name,
			styleStatus(component.Status),
		})
	}
	
	tbl := tablelist.NewModel(columns, rows, len(rows)+1) // +1 for header
	return HealthModel{tableModel: tbl}
}

func PrintHealthStatus(status *health.GetHealthOK) {
	m := NewHealthModel(status)
	fmt.Printf("Harbor Health Status:: %s\n", styleStatus(status.Payload.Status))

	p := tea.NewProgram(m.tableModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		return
	}
}