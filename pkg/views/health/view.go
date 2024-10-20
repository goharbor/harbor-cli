package health

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/health"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Component", Width: 18},
	{Title: "Status", Width: 26},
}

func PrintHealthStatus(status *health.GetHealthOK) {
	var rows []table.Row
	fmt.Printf("Harbor Health Status:: %s\n", styleStatus(status.Payload.Status))
	for _, component := range status.Payload.Components {
		rows = append(rows, table.Row{
			component.Name,
			styleStatus(component.Status),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func styleStatus(status string) string {
	if status == "healthy" {
		return views.GreenStyle.Render(status)
	}
	return views.RedStyle.Render(status)
}
