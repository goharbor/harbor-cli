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
package tablegridv2

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

type gridLoadedMsg struct {
	config Config
}

type gridLoadFailedMsg struct {
	err error
}

type CellStatus bool

type Styles struct {
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	Disabled   lipgloss.Style
	Header     lipgloss.Style
	Cursor     string
}

type Icons struct {
	Selected   string
	Unselected string
	Empty      string
}

type Config struct {
	RowLabels    []string
	ColLabels    []string
	Data         [][]CellStatus
	Disabled     map[int]map[int]bool
	ColumnWidths []int
	Styles       *Styles
	Icons        *Icons
	Footer       string
}

type Loader func() (Config, error)

type Model struct {
	Table       table.Model
	Data        [][]CellStatus
	RowLabels   []string
	ColLabels   []string
	Disabled    map[int]map[int]bool
	SelectedCol int
	Styles      Styles
	Icons       Icons
	Footer      string
	Load        Loader
	Spinner     spinner.Model
	Err         error
	LoadingText string
	done        bool
	animate     bool
	state       state
}

func NewModel(loadingText string, load Loader) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	if loadingText == "" {
		loadingText = "Loading..."
	}

	return Model{
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
	case gridLoadedMsg:
		m.applyConfig(msg.config)
		m.state = stateReady
		return m, nil
	case gridLoadFailedMsg:
		m.Err = msg.err
		m.state = stateError
		return m, tea.Quit
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+a":
			for rowIdx := range m.RowLabels {
				if m.Disabled != nil && m.Disabled[rowIdx] != nil {
					for colIdx := 1; colIdx < len(m.ColLabels); colIdx++ {
						if m.Disabled[rowIdx][colIdx] {
							continue
						}
						m.Data[rowIdx][colIdx-1] = true
					}
				} else {
					for colIdx := 1; colIdx < len(m.ColLabels); colIdx++ {
						m.Data[rowIdx][colIdx-1] = true
					}
				}
			}
			m.refreshTable(m.Table.Cursor(), m.SelectedCol)
			return m, nil
		case "ctrl+d":
			for rowIdx := range m.RowLabels {
				if m.Disabled != nil && m.Disabled[rowIdx] != nil {
					for colIdx := 1; colIdx < len(m.ColLabels); colIdx++ {
						if m.Disabled[rowIdx][colIdx] {
							continue
						}
						m.Data[rowIdx][colIdx-1] = false
					}
				} else {
					for colIdx := 1; colIdx < len(m.ColLabels); colIdx++ {
						m.Data[rowIdx][colIdx-1] = false
					}
				}
			}
			m.refreshTable(m.Table.Cursor(), m.SelectedCol)
			return m, nil
		case "ctrl+j":
			if m.Table.Cursor() < 0 || m.Table.Cursor() >= len(m.RowLabels) {
				return m, nil
			}
			rowIdx := m.Table.Cursor()
			for colIdx := 1; colIdx < len(m.ColLabels); colIdx++ {
				if m.Disabled != nil && m.Disabled[rowIdx] != nil && m.Disabled[rowIdx][colIdx] {
					continue
				}
				m.Data[rowIdx][colIdx-1] = true
			}
			m.refreshTable(rowIdx, m.SelectedCol)
			return m, nil
		case "ctrl+k":
			if m.Table.Cursor() < 0 || m.Table.Cursor() >= len(m.RowLabels) {
				return m, nil
			}
			rowIdx := m.Table.Cursor()
			for colIdx := 1; colIdx < len(m.ColLabels); colIdx++ {
				if m.Disabled != nil && m.Disabled[rowIdx] != nil && m.Disabled[rowIdx][colIdx] {
					continue
				}
				m.Data[rowIdx][colIdx-1] = false
			}
			m.refreshTable(rowIdx, m.SelectedCol)
			return m, nil
		case "ctrl+h":
			if m.SelectedCol < 1 || m.SelectedCol >= len(m.ColLabels) {
				return m, nil
			}
			colIdx := m.SelectedCol
			for rowIdx := range m.RowLabels {
				if m.Disabled != nil && m.Disabled[rowIdx] != nil && m.Disabled[rowIdx][colIdx] {
					continue
				}
				m.Data[rowIdx][colIdx-1] = true
			}
			m.refreshTable(m.Table.Cursor(), m.SelectedCol)
			return m, nil
		case "ctrl+l":
			if m.SelectedCol < 1 || m.SelectedCol >= len(m.ColLabels) {
				return m, nil
			}
			colIdx := m.SelectedCol
			for rowIdx := range m.RowLabels {
				if m.Disabled != nil && m.Disabled[rowIdx] != nil && m.Disabled[rowIdx][colIdx] {
					continue
				}
				m.Data[rowIdx][colIdx-1] = false
			}
			m.refreshTable(m.Table.Cursor(), m.SelectedCol)
			return m, nil
		case "ctrl+s":
			m.done = true
			return m, tea.Quit
		case "left", "h":
			curRow := m.Table.Cursor()
			for next := m.SelectedCol - 1; next >= 1; next-- {
				if m.Disabled == nil || m.Disabled[curRow] == nil || !m.Disabled[curRow][next] {
					m.SelectedCol = next
					m.refreshTable(curRow, m.SelectedCol)
					break
				}
			}
			return m, nil
		case "right", "l":
			curRow := m.Table.Cursor()
			for next := m.SelectedCol + 1; next < len(m.ColLabels); next++ {
				if m.Disabled == nil || m.Disabled[curRow] == nil || !m.Disabled[curRow][next] {
					m.SelectedCol = next
					m.refreshTable(curRow, m.SelectedCol)
					break
				}
			}
			return m, nil
		case "up", "k":
			m.Table, cmd = m.Table.Update(msg)
			for {
				r := m.Table.Cursor()
				if r <= 0 || m.Disabled == nil || m.Disabled[r] == nil || !m.Disabled[r][m.SelectedCol] {
					break
				}
				m.Table, _ = m.Table.Update(msg)
			}
			m.refreshTable(m.Table.Cursor(), m.SelectedCol)
			return m, cmd
		case "down", "j":
			m.Table, cmd = m.Table.Update(msg)
			for {
				r := m.Table.Cursor()
				if r >= len(m.RowLabels)-1 || m.Disabled == nil || m.Disabled[r] == nil || !m.Disabled[r][m.SelectedCol] {
					break
				}
				m.Table, _ = m.Table.Update(msg)
			}
			m.refreshTable(m.Table.Cursor(), m.SelectedCol)
			return m, cmd
		case "enter", " ":
			rowIdx := m.Table.Cursor()
			colIdx := m.SelectedCol
			if m.Disabled != nil && m.Disabled[rowIdx] != nil && m.Disabled[rowIdx][colIdx] {
				return m, nil
			}
			m.Data[rowIdx][colIdx-1] = !m.Data[rowIdx][colIdx-1]
			m.refreshTable(rowIdx, colIdx)
			return m, nil
		case "q", "ctrl+c":
			m.done = true
			return m, tea.Quit
		}
	}

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
		if m.done {
			return ""
		}
		cursor := m.Table.Cursor()
		m.refreshTable(cursor, m.SelectedCol)
		return m.Table.View() + m.footerText()
	}
}

func (m Model) loadCmd() tea.Cmd {
	return func() tea.Msg {
		config, err := m.Load()
		if err != nil {
			return gridLoadFailedMsg{err: err}
		}
		return gridLoadedMsg{config: config}
	}
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)
}

func DefaultStyles() Styles {
	return Styles{
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("42")),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		Disabled:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Header:     lipgloss.NewStyle().Bold(true),
		Cursor:     "▶",
	}
}

func DefaultIcons() Icons {
	return Icons{
		Selected:   "✅",
		Unselected: "❌",
		Empty:      " ",
	}
}

func (m *Model) applyConfig(config Config) {
	styles := DefaultStyles()
	if config.Styles != nil {
		styles = *config.Styles
	}

	icons := DefaultIcons()
	if config.Icons != nil {
		icons = *config.Icons
	}

	colWidths := config.ColumnWidths
	if colWidths == nil {
		colWidths = make([]int, len(config.ColLabels))
		for i := range colWidths {
			colWidths[i] = 16
			if i == 0 {
				colWidths[i] = 20
			}
		}
	}

	columns := make([]table.Column, len(config.ColLabels))
	for i, label := range config.ColLabels {
		columns[i] = table.Column{Title: label, Width: colWidths[i]}
	}

	data := config.Data
	if data == nil {
		data = make([][]CellStatus, len(config.RowLabels))
		for i := range data {
			data[i] = make([]CellStatus, len(config.ColLabels)-1)
		}
	}

	rows := buildRows(config.RowLabels, data, -1, -1, config.Disabled, styles, icons)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.Inherit(styles.Header)
	t.SetStyles(tableStyles)

	m.Table = t
	m.Data = data
	m.RowLabels = config.RowLabels
	m.ColLabels = config.ColLabels
	m.Disabled = config.Disabled
	m.SelectedCol = 1
	m.Styles = styles
	m.Icons = icons
	m.Footer = config.Footer
}

func buildRows(labels []string, data [][]CellStatus, highlightRow, highlightCol int, disabled map[int]map[int]bool, styles Styles, icons Icons) []table.Row {
	rows := make([]table.Row, len(labels))
	for i, label := range labels {
		cells := make([]string, len(data[i])+1)
		cells[0] = label
		for j := 0; j < len(data[i]); j++ {
			colIdx := j + 1
			if disabled != nil && disabled[i] != nil && disabled[i][colIdx] {
				cells[colIdx] = styles.Disabled.Render(icons.Empty)
				continue
			}
			var icon string
			if data[i][j] {
				icon = styles.Selected.Render(icons.Selected)
			} else {
				icon = styles.Unselected.Render(icons.Unselected)
			}
			if i == highlightRow && colIdx == highlightCol {
				cells[colIdx] = fmt.Sprintf("%s %s", styles.Cursor, icon)
			} else {
				cells[colIdx] = icon
			}
		}
		rows[i] = table.Row(cells)
	}
	return rows
}

func (m *Model) refreshTable(highlightRow, highlightCol int) {
	m.Table.SetRows(buildRows(m.RowLabels, m.Data, highlightRow, highlightCol, m.Disabled, m.Styles, m.Icons))
}

func (m Model) footerText() string {
	if m.Footer != "" {
		return m.Footer
	}
	return "\n ↑/↓ move row • ⌃J toggle row on  • ⌃H toggle col on  • ^A toggle table on  • space/enter to toggle\n" +
		" ←/→ move col • ⌃K toggle row off • ⌃L toggle col off • ^D toggle table off • ^S submit • q to cancel \n"
}
