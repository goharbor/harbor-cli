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
package tablelist

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/harbor-cli/pkg/views"
)

const (
	WidthXS  = 4
	WidthS   = 8
	WidthM   = 12
	WidthL   = 16
	WidthXL  = 20
	WidthXXL = 24
	Width3XL = 30
)

type Model struct {
	Table     table.Model
	Spinner   spinner.Model
	IsLoading bool
}

func NewModel(columns []table.Column, rows []table.Row, height int) Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height+1),
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

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		Table:     t,
		Spinner:   sp,
		IsLoading: len(rows) == 0,
	}
}

func (m Model) Init() tea.Cmd {
	if m.IsLoading {
		return m.Spinner.Tick
	}
	return tea.Quit
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.IsLoading {
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}
	case []table.Row:
		m.Table.SetRows(msg)
		m.IsLoading = false
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.IsLoading {
		return "\n  " + m.Spinner.View() + " Loading data...\n"
	}
	return views.BaseStyle.Render(m.Table.View()) + "\n"
}
