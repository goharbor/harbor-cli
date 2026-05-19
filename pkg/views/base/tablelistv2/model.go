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
	"time"

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

type DataLoadedMessage struct {
	Data []table.Row
	Err  error
}

type Model struct {
	Width     int // Width is required since the v2 fails with no explicit width defined
	Table     table.Model
	Error     error
	Columns   []table.Column
	IsLoading bool
	Spinner   spinner.Model
	LoadData  func() ([]table.Row, error)
}

func NewModel(columns []table.Column, fn func() ([]table.Row, error)) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	w := 0
	for _, v := range columns {
		w += v.Width
	}

	return Model{
		LoadData:  fn,
		Columns:   columns,
		IsLoading: true,
		Spinner:   s,
		Width:     w,
	}
}

func (m Model) Init() tea.Cmd {
	// Calls the Data function asynchronously and also starts ticking
	// Load for animation
	return tea.Batch(tea.Tick(0, func(_ time.Time) tea.Msg {
		data, err := m.LoadData()
		return DataLoadedMessage{Data: data, Err: err}
	}),
		m.Spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case DataLoadedMessage:
		if msg.Err != nil {
			m.Error = msg.Err
			return m, tea.Quit
		}
		t := table.New(
			table.WithColumns(m.Columns),
			table.WithRows(msg.Data),
			table.WithFocused(true),
			table.WithWidth(m.Width),
			table.WithHeight(len(msg.Data)+1),
		)
		t.SetStyles(tableStyles())

		m.Table = t
		m.IsLoading = false

		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.IsLoading {
		return fmt.Sprintf("%s Loading...", m.Spinner.View())
	}

	return views.BaseStyle.Render(m.Table.View()) + "\n"
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
