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

package loadingtable

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

type FetchMsg struct {
	Rows []table.Row
	Err  error
}

type Model struct {
	spinner  spinner.Model
	table    tablelist.Model
	columns  []table.Column
	fetcher  tea.Cmd
	title    string
	loading  bool
	err      error
	Aborted  bool
	Choice   table.Row
	Selected bool
}

func NewModel(title string, columns []table.Column, fetcher tea.Cmd) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		spinner: s,
		columns: columns,
		fetcher: fetcher,
		title:   title,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetcher)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.Aborted = true
			return m, tea.Quit
		case "enter":
			if !m.loading && m.err == nil {
				m.Choice = m.table.Table.SelectedRow()
				m.Selected = true
				return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case FetchMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, tea.Quit
		}
		m.table = tablelist.NewModel(m.columns, msg.Rows, len(msg.Rows))
		return m, nil
	}

	if !m.loading && m.err == nil {
		var cmd tea.Cmd
		m.table.Table, cmd = m.table.Table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	titleStr := ""
	if m.title != "" {
		titleStr = views.TitleStyle.Render(m.title) + "\n"
	}

	if m.err != nil {
		return titleStr + views.BaseStyle.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Error: %v", m.err))) + "\n"
	}
	if m.loading {
		return titleStr + views.BaseStyle.Render(fmt.Sprintf("%s Loading data...", m.spinner.View())) + "\n"
	}
	return titleStr + m.table.View()
}
