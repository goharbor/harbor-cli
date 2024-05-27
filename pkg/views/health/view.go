package health

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/health"
	"github.com/goharbor/harbor-cli/pkg/views"
)

type model struct {
	table table.Model
}

var columns = []table.Column{
	{Title: "Component Name", Width: 14},
	{Title: "Status", Width: 20},
	{Title: "Error Message", Width: 32},
}

func (m model) Init() tea.Cmd {
	return tea.Quit
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return views.BaseStyle.Render(m.table.View()) + "\n"
}

func View(health *health.GetHealthOK) {
	style := views.TitleStyle.Bold(true)
	fmt.Print(style.Render("Status: "))
	statusStyle := getStatusStyle(health.Payload.Status)
	fmt.Println(statusStyle.Render(health.Payload.Status))

	var rows []table.Row
	for _, component := range health.Payload.Components {
		statusStyle := getStatusStyle(component.Status)
		status := statusStyle.Render(component.Status)
		rows = append(rows, table.Row{
			component.Name,
			status,
			component.Error,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
	)

	// Set the styles for the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(s.Cell.GetForeground()).
		Background(s.Cell.GetBackground()).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getStatusStyle(status string) lipgloss.Style {
	statusStyle := views.RedStyle
	if status == "healthy" {
		statusStyle = views.GreenStyle
	}
	return statusStyle
}
