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
package tablegrid

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyMap defines all key bindings for TableGrid and implements help.KeyMap.
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Toggle   key.Binding
	RowOn    key.Binding
	RowOff   key.Binding
	ColOn    key.Binding
	ColOff   key.Binding
	TableOn  key.Binding
	TableOff key.Binding
	Submit   key.Binding
	Quit     key.Binding
	Help     key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Left, k.Right,
		k.Toggle, k.TableOn,
		k.Submit, k.Quit, k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Toggle, k.RowOn, k.RowOff},
		{k.ColOn, k.ColOff, k.TableOn, k.TableOff},
		{k.Submit, k.Quit},
	}
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
		Left:     key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "move left")),
		Right:    key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "move right")),
		Toggle:   key.NewBinding(key.WithKeys(" ", "enter"), key.WithHelp("space/enter", "toggle cell")),
		RowOn:    key.NewBinding(key.WithKeys("ctrl+j"), key.WithHelp("^J", "row on")),
		RowOff:   key.NewBinding(key.WithKeys("ctrl+k"), key.WithHelp("^K", "row off")),
		ColOn:    key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp("^H", "col on")),
		ColOff:   key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("^L", "col off")),
		TableOn:  key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("^A", "all on")),
		TableOff: key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("^D", "all off")),
		Submit:   key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("^S", "submit")),
		Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "See For How to Toggle")),
	}
}

// CellStatus represents a cell's toggle state
type CellStatus bool

// TableGrid is a generic interactive grid for toggling options
type TableGrid struct {
	Table       table.Model
	Data        [][]CellStatus       // 2D grid of cell states (true = selected, false = not selected)
	RowLabels   []string             // Labels for rows
	ColLabels   []string             // Labels for columns
	Disabled    map[int]map[int]bool // Which cells are disabled
	SelectedCol int                  // Currently selected column
	Styles      Styles               // Custom styles
	Icons       Icons                // Custom icons
	Footer      string               // Custom footer text
	Help        help.Model
	Keys        KeyMap
}

// Styles contains customizable styles for the table grid
type Styles struct {
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	Disabled   lipgloss.Style
	Header     lipgloss.Style
	Cursor     string // Cursor indicator
}

// Icons defines how cells are displayed
type Icons struct {
	Selected   string // Icon for selected cells
	Unselected string // Icon for unselected cells
	Empty      string // Icon for disabled cells
}

// Config holds parameters for creating a new TableGrid
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

// DefaultStyles returns the default styles
func DefaultStyles() Styles {
	return Styles{
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("42")),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		Disabled:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Header:     lipgloss.NewStyle().Bold(true),
		Cursor:     "▶",
	}
}

// DefaultIcons returns the default icons
func DefaultIcons() Icons {
	return Icons{
		Selected:   "✅",
		Unselected: "❌",
		Empty:      " ",
	}
}

// New creates a new TableGrid
func New(config Config) *TableGrid {
	// Apply defaults
	styles := DefaultStyles()
	if config.Styles != nil {
		styles = *config.Styles
	}

	icons := DefaultIcons()
	if config.Icons != nil {
		icons = *config.Icons
	}

	// Set column widths
	colWidths := config.ColumnWidths
	if colWidths == nil {
		colWidths = make([]int, len(config.ColLabels))
		for i := range colWidths {
			colWidths[i] = 16 // Default width
			if i == 0 {
				colWidths[i] = 20 // Wider for first column
			}
		}
	}

	// Create columns
	columns := make([]table.Column, len(config.ColLabels))
	for i, label := range config.ColLabels {
		columns[i] = table.Column{
			Title: label,
			Width: colWidths[i],
		}
	}

	// Initialize data grid if not provided
	data := config.Data
	if data == nil {
		data = make([][]CellStatus, len(config.RowLabels))
		for i := range data {
			data[i] = make([]CellStatus, len(config.ColLabels)-1)
		}
	}

	// Build initial rows
	rows := buildRows(config.RowLabels, data, -1, -1, config.Disabled, styles, icons)

	// Create table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)

	// Apply table styles
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.Inherit(styles.Header)
	t.SetStyles(tableStyles)

	return &TableGrid{
		Table:       t,
		Data:        data,
		RowLabels:   config.RowLabels,
		ColLabels:   config.ColLabels,
		Disabled:    config.Disabled,
		SelectedCol: 1, // Start with first editable column
		Styles:      styles,
		Icons:       icons,
		Footer:      config.Footer,
		Help:        help.New(),
		Keys:        DefaultKeyMap(),
	}
}

// buildRows constructs the table rows with proper styling
func buildRows(
	labels []string,
	data [][]CellStatus,
	highlightRow, highlightCol int,
	disabled map[int]map[int]bool,
	styles Styles,
	icons Icons,
) []table.Row {
	rows := make([]table.Row, len(labels))

	for i, label := range labels {
		cells := make([]string, len(data[i])+1) // +1 for label column
		cells[0] = label

		for j := 0; j < len(data[i]); j++ {
			colIdx := j + 1 // Adjust for label column

			// Handle disabled cells
			if disabled != nil && disabled[i] != nil && disabled[i][colIdx] {
				cells[colIdx] = styles.Disabled.Render(icons.Empty)
				continue
			}

			// Render cell with appropriate icon
			var icon string
			if data[i][j] {
				icon = styles.Selected.Render(icons.Selected)
			} else {
				icon = styles.Unselected.Render(icons.Unselected)
			}

			// Add cursor if cell is highlighted
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

// Init initializes the model
func (m *TableGrid) Init() tea.Cmd {
	return nil
}

// Update handles UI updates
func (m *TableGrid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+a":
			// Turn all cells on
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
			// Turn all cells off
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
			// Turn on all in current row
			if m.Table.Cursor() < 0 || m.Table.Cursor() >= len(m.RowLabels) {
				return m, nil // No valid row selected
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
			// Turn off all in current row
			if m.Table.Cursor() < 0 || m.Table.Cursor() >= len(m.RowLabels) {
				return m, nil // No valid row selected
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
			// Turn on all in current column
			if m.SelectedCol < 1 || m.SelectedCol >= len(m.ColLabels) {
				return m, nil // No valid column selected
			}
			colIdx := m.SelectedCol
			for rowIdx := 0; rowIdx < len(m.RowLabels); rowIdx++ {
				if m.Disabled != nil && m.Disabled[rowIdx] != nil && m.Disabled[rowIdx][colIdx] {
					continue
				}
				m.Data[rowIdx][colIdx-1] = true
			}
			m.refreshTable(m.Table.Cursor(), m.SelectedCol)
			return m, nil

		case "ctrl+l":
			// Turn off all in current column
			if m.SelectedCol < 1 || m.SelectedCol >= len(m.ColLabels) {
				return m, nil // No valid column selected
			}
			colIdx := m.SelectedCol
			for rowIdx := 0; rowIdx < len(m.RowLabels); rowIdx++ {
				if m.Disabled != nil && m.Disabled[rowIdx] != nil && m.Disabled[rowIdx][colIdx] {
					continue
				}
				m.Data[rowIdx][colIdx-1] = false
			}
			m.refreshTable(m.Table.Cursor(), m.SelectedCol)
			return m, nil

		case "ctrl+s":
			return m, tea.Quit

		case "left", "h":
			// Move selection left
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
			// Move selection right
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
			// Toggle cell
			rowIdx := m.Table.Cursor()
			colIdx := m.SelectedCol
			if m.Disabled != nil && m.Disabled[rowIdx] != nil && m.Disabled[rowIdx][colIdx] {
				return m, nil
			}
			m.Data[rowIdx][colIdx-1] = !m.Data[rowIdx][colIdx-1]
			m.refreshTable(rowIdx, colIdx)
			return m, nil

		case "?":
			m.Help.ShowAll = !m.Help.ShowAll
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

// refreshTable updates the table with current data and highlights
func (m *TableGrid) refreshTable(highlightRow, highlightCol int) {
	newRows := buildRows(m.RowLabels, m.Data, highlightRow, highlightCol, m.Disabled, m.Styles, m.Icons)
	m.Table.SetRows(newRows)
}

// View renders the component
func (m *TableGrid) View() string {
	cursor := m.Table.Cursor()
	m.refreshTable(cursor, m.SelectedCol)
	out := m.Table.View()

	return out + "\n" + m.Help.View(m.Keys)
}

// GetData returns the current selection state
func (m *TableGrid) GetData() [][]CellStatus {
	return m.Data
}

// GetSelectedItems returns the coordinates of all selected cells
func (m *TableGrid) GetSelectedItems() [][]int {
	var selected [][]int

	for i := range m.Data {
		for j := range m.Data[i] {
			if m.Data[i][j] {
				selected = append(selected, []int{i, j})
			}
		}
	}

	return selected
}
