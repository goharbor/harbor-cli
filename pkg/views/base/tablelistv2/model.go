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
package tablelistv2

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/harbor-cli/pkg/views"
	"golang.org/x/term"
)

type state string

const (
	stateLoading state = "loading"
	stateReady   state = "ready"
	stateError   state = "error"
)

type dataLoadedMsg struct {
	rows []table.Row
}

type dataLoadFailedMsg struct {
	err error
}

type Loader func() ([]table.Row, error)

type Model struct {
	Table       table.Model
	Columns     []table.Column
	Load        Loader
	Spinner     spinner.Model
	Err         error
	LoadingText string
	animate     bool
	state       state
}

func NewModel(columns []table.Column, loadingText string, load Loader) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	if loadingText == "" {
		loadingText = "Loading..."
	}

	return Model{
		Columns:     columns,
		Load:        load,
		Spinner:     s,
		LoadingText: loadingText,
		animate:     term.IsTerminal(int(os.Stdout.Fd())), // #nosec G115 - fd fits in int on all supported platforms
		state:       stateLoading,
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.loadCmd()}
	if m.animate {
		cmds = append(cmds, m.Spinner.Tick)
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.state != stateLoading || !m.animate {
			return m, nil
		}

		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case dataLoadedMsg:
		m.Table = table.New(
			table.WithColumns(m.Columns),
			table.WithRows(msg.rows),
			table.WithFocused(true),
			table.WithHeight(len(msg.rows)+1),
			table.WithWidth(columnsWidth(m.Columns)),
		)
		m.Table.SetStyles(tableStyles())
		m.state = stateReady
		return m, tea.Quit
	case dataLoadFailedMsg:
		m.Err = msg.err
		m.state = stateError
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		if m.animate {
			return views.BaseStyle.Render(fmt.Sprintf("%s %s", m.Spinner.View(), m.LoadingText)) + "\n"
		}
		return views.BaseStyle.Render(m.LoadingText) + "\n"
	case stateError:
		return views.BaseStyle.Render(errorStyle().Render(m.Err.Error())) + "\n"
	default:
		return views.BaseStyle.Render(m.Table.View()) + "\n"
	}
}

func (m Model) loadCmd() tea.Cmd {
	return func() tea.Msg {
		rows, err := m.Load()
		if err != nil {
			return dataLoadFailedMsg{err: err}
		}
		return dataLoadedMsg{rows: rows}
	}
}

func columnsWidth(columns []table.Column) int {
	width := 0
	for _, column := range columns {
		width += column.Width
	}
	return width
}

func tableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(s.Cell.GetForeground()).
		Background(s.Cell.GetBackground()).
		Bold(false)

	return s
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)
}
