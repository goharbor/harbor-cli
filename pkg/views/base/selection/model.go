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
package selection

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/goharbor/harbor-cli/pkg/views"
)

const listHeight = 14

type Item string

func (i Item) FilterValue() string { return string(i) }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := views.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return views.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	List    list.Model
	Choice  string
	Aborted bool

	loading  bool
	err      error
	spinner  spinner.Model
	fetchCmd tea.Cmd
}

type DataLoadedMsg struct {
	Items []list.Item
	Err   error
}

type FetchFn func() ([]list.Item, error)

func FetchCmd(fn FetchFn) tea.Cmd {
	return func() tea.Msg {
		items, err := fn()
		return DataLoadedMsg{Items: items, Err: err}
	}
}

func NewModel(items []list.Item, construct string) Model {
	const defaultWidth = 20
	l := list.New(items, ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select a " + construct
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = views.TitleStyle
	l.Styles.PaginationStyle = views.PaginationStyle
	l.Styles.HelpStyle = views.HelpStyle

	return Model{
		List:    l,
		loading: false,
	}
}

func NewModelWithFetch(fetchFn FetchFn, construct string) Model {
	m := NewModel(nil, construct)
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
	switch msg := msg.(type) {
	case DataLoadedMsg:
		m.loading = false
		m.err = msg.Err
		if msg.Err == nil {
			m.List.SetItems(msg.Items)
		}
		return m, nil
	case spinner.TickMsg:
		if m.loading {
			var spinnerCmd tea.Cmd
			m.spinner, spinnerCmd = m.spinner.Update(msg)
			return m, spinnerCmd
		}
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.List.FilterState() != list.Filtering {
				i, ok := m.List.SelectedItem().(Item)
				if ok {
					m.Choice = string(i)
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	if m.loading {
		return tea.NewView(m.spinner.View() + " Loading selections...\n")
	}
	if m.err != nil {
		return tea.NewView("Error: " + m.err.Error() + "\n")
	}
	if m.Choice != "" {
		return tea.NewView("")
	}
	return tea.NewView("\n" + m.List.View())
}
