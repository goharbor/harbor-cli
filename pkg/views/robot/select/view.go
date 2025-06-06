// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package robot

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

const (
	WidthResource = 20
	WidthAction   = 16
)

var columns = []table.Column{
	{Title: "Resource", Width: WidthResource},
	{Title: "Create", Width: WidthAction},
	{Title: "Delete", Width: WidthAction},
	{Title: "List", Width: WidthAction},
	{Title: "Pull", Width: WidthAction},
	{Title: "Push", Width: WidthAction},
	{Title: "Read", Width: WidthAction},
	{Title: "Stop", Width: WidthAction},
	{Title: "Update", Width: WidthAction},
}

var resourceNames = []string{
	"Accessory",
	"Artifact",
	"Artifact Addition",
	"Artifact Label",
	"Export CVE",
	"Immutable Tag",
	"Label",
	"Log",
	"Member",
	"Metadata",
	"Notification Policy",
	"Preheat Policy",
	"Project",
	"Quota",
	"Repository",
	"Robot Account",
	"SBOM",
	"Scan",
	"Scanner",
	"Tag",
	"Tag Retention",
}

type Model struct {
	Table    table.Model
	perms    [][]bool
	selCol   int
	disabled map[int]map[int]bool
}

// ToDo: Move this to package base/tablegrid
func NewModel() Model {
	numRows := len(resourceNames)
	numCols := len(columns) - 1
	perms := make([][]bool, numRows)
	for i := 0; i < numRows; i++ {
		perms[i] = make([]bool, numCols)
	}
	disabled := map[int]map[int]bool{
		0:  {1: true, 2: true, 3: false, 4: true, 5: true, 6: true, 7: true, 8: true},
		1:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: true},
		2:  {1: true, 2: true, 3: true, 4: true, 5: true, 6: false, 7: true, 8: true},
		3:  {1: false, 2: false, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true},
		4:  {1: false, 2: true, 3: true, 4: true, 5: true, 6: false, 7: true, 8: true},
		5:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: true, 7: true, 8: false},
		6:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		7:  {1: true, 2: true, 3: false, 4: true, 5: true, 6: true, 7: true, 8: true},
		8:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		9:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		10: {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		11: {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		12: {1: true, 2: true, 3: true, 4: true, 5: true, 6: false, 7: false, 8: true},
		13: {1: true, 2: true, 3: true, 4: true, 5: true, 6: false, 7: true, 8: true},
		14: {1: true, 2: false, 3: false, 4: false, 5: false, 6: false, 7: true, 8: false},
		15: {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: true},
		16: {1: false, 2: true, 3: true, 4: true, 5: true, 6: false, 7: false, 8: true},
		17: {1: false, 2: true, 3: true, 4: true, 5: true, 6: false, 7: false, 8: true},
		18: {1: false, 2: true, 3: true, 4: true, 5: true, 6: false, 7: true, 8: true},
		19: {1: false, 2: false, 3: false, 4: true, 5: true, 6: true, 7: true, 8: true},
		20: {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
	}
	rows := buildRows(resourceNames, perms, -1, -1, disabled)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)
	styles := table.DefaultStyles()
	styles.Header = styles.Header.Bold(false)
	t.SetStyles(styles)
	return Model{
		Table:    t,
		perms:    perms,
		selCol:   1,
		disabled: disabled,
	}
}

func buildRows(
	names []string,
	perms [][]bool,
	highlightRow, highlightCol int,
	disabled map[int]map[int]bool,
) []table.Row {
	grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	rows := make([]table.Row, len(names))
	for i, name := range names {
		cells := make([]string, len(columns))
		cells[0] = name
		for j := 1; j < len(columns); j++ {
			if disabled[i] != nil && disabled[i][j] {
				cells[j] = grayStyle.Render(" ")
				continue
			}
			var icon string
			if perms[i][j-1] {
				icon = "✅"
			} else {
				icon = "❌"
			}
			if i == highlightRow && j == highlightCol {
				cells[j] = fmt.Sprintf("▶ %s", icon)
			} else {
				cells[j] = icon
			}
		}
		rows[i] = table.Row(cells)
	}
	return rows
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+a":
			rowIdx := m.Table.Cursor()
			for colIdx := 1; colIdx < len(columns); colIdx++ {
				if m.disabled[rowIdx] != nil && m.disabled[rowIdx][colIdx] {
					continue
				}
				m.perms[rowIdx][colIdx-1] = !m.perms[rowIdx][colIdx-1]
			}
			newRows := buildRows(resourceNames, m.perms, rowIdx, m.selCol, m.disabled)
			m.Table.SetRows(newRows)
			return m, nil

		case "ctrl+s":
			return m, tea.Quit

		case "left", "h":
			curRow := m.Table.Cursor()
			for next := m.selCol - 1; next >= 1; next-- {
				if m.disabled[curRow] == nil || !m.disabled[curRow][next] {
					m.selCol = next
					break
				}
			}
			return m, nil

		case "right", "l":
			curRow := m.Table.Cursor()
			for next := m.selCol + 1; next < len(columns); next++ {
				if m.disabled[curRow] == nil || !m.disabled[curRow][next] {
					m.selCol = next
					break
				}
			}
			return m, nil

		case "up", "k":
			m.Table, cmd = m.Table.Update(msg)
			for {
				r := m.Table.Cursor()
				if r <= 0 || m.disabled[r] == nil || !m.disabled[r][m.selCol] {
					break
				}
				m.Table, _ = m.Table.Update(msg)
			}
			return m, cmd

		case "down", "j":
			m.Table, cmd = m.Table.Update(msg)
			for {
				r := m.Table.Cursor()
				if r >= len(resourceNames)-1 || m.disabled[r] == nil || !m.disabled[r][m.selCol] {
					break
				}
				m.Table, _ = m.Table.Update(msg)
			}
			return m, cmd

		case "enter", " ":
			rowIdx := m.Table.Cursor()
			colIdx := m.selCol
			if m.disabled[rowIdx] != nil && m.disabled[rowIdx][colIdx] {
				return m, nil
			}
			m.perms[rowIdx][colIdx-1] = !m.perms[rowIdx][colIdx-1]
			newRows := buildRows(resourceNames, m.perms, rowIdx, colIdx, m.disabled)
			m.Table.SetRows(newRows)
			return m, nil

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	cursor := m.Table.Cursor()
	highlighted := buildRows(resourceNames, m.perms, cursor, m.selCol, m.disabled)
	orig := m.Table.Rows()
	m.Table.SetRows(highlighted)
	out := m.Table.View()
	m.Table.SetRows(orig)
	footer := "\n ↑/↓ move row • ←/→ move col • space/enter to toggle-and-submit • ⌃A toggle row • q to cancel\n"
	return out + footer
}

func ListPermissions(perms *models.Permissions, ch chan<- []models.Permission) {
	t := NewModel()
	final, err := tea.NewProgram(t, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	m := final.(Model)
	selectedPerms := []models.Permission{}
	for rowIdx, res := range resourceNames {
		kebabResource := strings.ToLower(res)
		kebabResource = strings.ReplaceAll(kebabResource, " ", "-")
		for colIdx := 1; colIdx < len(columns); colIdx++ {
			if m.perms[rowIdx][colIdx-1] {
				action := strings.ToLower(columns[colIdx].Title)
				selectedPerms = append(selectedPerms, models.Permission{
					Resource: kebabResource,
					Action:   action,
				})
			}
		}
	}
	ch <- selectedPerms
}

func ListRobot(robots []*models.Robot, choice chan<- int64) {
	itemsList := make([]list.Item, len(robots))
	items := map[string]int64{}
	for i, r := range robots {
		items[r.Name] = r.ID
		itemsList[i] = selection.Item(r.Name)
	}
	m := selection.NewModel(itemsList, "Robot")
	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if pm, ok := p.(selection.Model); ok {
		choice <- items[pm.Choice]
	}
}
