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
package listselect

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

func (i Item) FilterValue() string { return string(i) }

type ItemDelegate struct {
	Selected *map[int]struct{}
}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	checked := " "
	if d.Selected != nil {
		if _, ok := (*d.Selected)[index]; ok {
			checked = "✓"
		}
	}

	str := fmt.Sprintf("[%s] %d. %s", checked, index+1, i)

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
	Choices  []string
	Selected map[int]struct{}
	Aborted  bool
}

func NewModel(items []list.Item, construct string) Model {
	const defaultWidth = 20
	selected := make(map[int]struct{})
	l := list.New(items, ItemDelegate{Selected: &selected}, defaultWidth, listHeight)
	l.Title = "Select one or more " + construct + " (space to toggle, enter to confirm)"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = views.TitleStyle
	l.Styles.PaginationStyle = views.PaginationStyle
	l.Styles.HelpStyle = views.HelpStyle

	return Model{List: l, Selected: selected}
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
		case " ":
			idx := m.List.Index()
			if _, ok := m.Selected[idx]; ok {
				delete(m.Selected, idx)
			} else {
				m.Selected[idx] = struct{}{}
			}
			return m, nil
		case "enter":
			if len(m.Selected) == 0 {
				cmd := m.List.NewStatusMessage("!! Please select at least one item !!")
				return m, cmd
			}
			for idx := range m.Selected {
				if i, ok := m.List.Items()[idx].(Item); ok {
					m.Choices = append(m.Choices, string(i))
				}
			}
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.Aborted = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Aborted {
		return ""
	}
	if len(m.Choices) > 0 {
		return fmt.Sprintf("Selected: %s\n", strings.Join(m.Choices, ", "))
	}
	return "\n" + m.List.View()
}
