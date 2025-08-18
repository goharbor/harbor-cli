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
package multiselect

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/views"
)

const listHeight = 14

type Item string

func (i Item) FilterValue() string { return "" }

type ItemDelegate struct {
	Selected *map[int]bool // Use pointer to share the same map
}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	checkbox := "☐"
	if d.Selected != nil && (*d.Selected)[index] {
	if d.Selected != nil {
		if selected, ok := (*d.Selected)[index]; ok && selected {
			checkbox = "☒"
		}
	}

	str := fmt.Sprintf("%s %d. %s", checkbox, index+1, i)

	fn := views.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return views.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	List     list.Model
	Choice   []string
	Selected map[int]bool
	Aborted  bool
}

func NewModel(items []list.Item, construct string) Model {
	const defaultWidth = 20

	selected := make(map[int]bool)
	delegate := ItemDelegate{Selected: &selected}

	l := list.New(items, delegate, defaultWidth, listHeight)
	l.Title = "Select " + construct + " (Space to toggle, Enter to confirm)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = views.TitleStyle
	l.Styles.PaginationStyle = views.PaginationStyle
	l.Styles.HelpStyle = views.HelpStyle

	return Model{
		List:     l,
		Selected: selected,
		Choice:   make([]string, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.Aborted = true
			return m, tea.Quit

		case " ":
			currentIndex := m.List.Index()
			m.Selected[currentIndex] = !m.Selected[currentIndex]
			return m, nil

		case "enter":
			m.Choice = make([]string, 0)
			for i, isSelected := range m.Selected {
				if isSelected && i < len(m.List.Items()) {
			for i, listItem := range m.List.Items() {
				if m.Selected[i] {
					if item, ok := listItem.(Item); ok {
						m.Choice = append(m.Choice, string(item))
					}
				}
			}

			if len(m.Choice) == 0 {
				return m, nil
			}

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if len(m.Choice) > 0 {
		return ""
	}

	selectedCount := 0
	for _, selected := range m.Selected {
		if selected {
			selectedCount++
		}
	}

	helpText := fmt.Sprintf("\nSelected: %d | Space: toggle, Enter: confirm, q: quit", selectedCount)

	return "\n" + m.List.View() + helpText
}
