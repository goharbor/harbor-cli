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
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	Table    table.Model
	loading  bool
	err      error
	spinner  spinner.Model
	fetchCmd tea.Cmd
}

type DataLoadedMsg struct {
	Rows []table.Row
	Err  error
}

type FetchFn func() ([]table.Row, error)

func FetchCmd(fn FetchFn) tea.Cmd {
	return func() tea.Msg {
		rows, err := fn()
		return DataLoadedMsg{Rows: rows, Err: err}
	}
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

	return Model{
		Table:   t,
		loading: false,
	}
}

func NewModelWithFetch(columns []table.Column, fetchFn FetchFn, height int) Model {
	m := NewModel(columns, nil, height)
	m.loading = true
	m.fetchCmd = FetchCmd(fetchFn)
	
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m.spinner = s
	
	return m
}

func (m Model) Init() tea.Cmd {
	if m.loading {
		return tea.Batch(m.fetchCmd, m.spinner.Tick)
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case DataLoadedMsg:
		m.loading = false
		m.err = msg.Err
		if msg.Err == nil {
			m.Table.SetRows(msg.Rows)
		}
		return m, nil
	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	if m.loading {
		return tea.NewView(m.spinner.View() + " Loading data...\n")
	}
	if m.err != nil {
		return tea.NewView("Error: " + m.err.Error() + "\n")
	}
	return tea.NewView(views.BaseStyle.Render(m.Table.View()) + "\n")
}
