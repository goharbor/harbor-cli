package list

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views"
)

type model struct {
	table table.Model
}

var columns = []table.Column{
	{Title: "ID", Width: 4},
	{Title: "Member Name", Width: 12},
	{Title: "Project ID", Width: 12},
	{Title: "Role", Width: 12},
}

var columnsWide = []table.Column{
	{Title: "ID", Width: 4},
	{Title: "Member Name", Width: 12},
	{Title: "Type", Width: 8},
	{Title: "Project ID", Width: 12},
	{Title: "Role ID", Width: 8},
	{Title: "Role Name", Width: 12},
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

func ListMembers(members []*models.ProjectMemberEntity, wide bool) {
	var rows []table.Row
	for _, member := range members {
		memberID := strconv.FormatInt(member.ID, 10)
		projectID := strconv.FormatInt(member.ProjectID, 10)
		// roleName := utils.CamelCaseToHR(member.RoleName)
		roleName := member.RoleName

		if wide {
			roleID := strconv.FormatInt(member.RoleID, 10)
			memberType := member.EntityType

			if memberType == "u" {
				memberType = "User"
			}

			rows = append(rows, table.Row{
				memberID,
				member.EntityName,
				memberType,
				projectID,
				roleID,
				roleName,
			})
		} else {
			rows = append(rows, table.Row{
				memberID, // Member Name
				member.EntityName,
				projectID,
				roleName,
			})
		}
	}

	cols := columns
	if wide {
		cols = columnsWide
	}

	t := table.New(
		table.WithRows(rows),
		table.WithColumns(cols),
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
