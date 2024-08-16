package health

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/health"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Component", Width: 20},
	{Title: "Status", Width: 20},
}

type HealthModel struct {
	tableModel tablelist.Model
}

func NewHealthModel(status *health.GetHealthOK) HealthModel {
	var rows []table.Row

	rows = append(rows, table.Row{"Overall", status.Payload.Status})

	for _, component := range status.Payload.Components {
		rows = append(rows, table.Row{
			component.Name,
			component.Status,
		})
	}
	
	tbl := tablelist.NewModel(columns, rows, len(rows)+1) // +1 for header
	return HealthModel{tableModel: tbl}
}

func PrintHealthStatus(status *health.GetHealthOK) {
	m := NewHealthModel(status)
	fmt.Println("Harbor Health Status:")
	fmt.Printf("Status: %s\n", status.Payload.Status)

	p := tea.NewProgram(m.tableModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		return
	}
}