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
package selectionv2

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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

type itemsLoadedMsg struct {
	items []list.Item
}

type itemsLoadFailedMsg struct {
	err error
}

type Item string

func (i Item) FilterValue() string { return string(i) }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, item)

	fn := views.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return views.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Loader func() ([]list.Item, error)

type Model struct {
	List        list.Model
	Load        Loader
	Spinner     spinner.Model
	Choice      string
	Aborted     bool
	done        bool
	Err         error
	Construct   string
	LoadingText string
	animate     bool
	state       state
}

func NewModel(construct, loadingText string, load Loader) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	if loadingText == "" {
		loadingText = "Loading..."
	}

	return Model{
		Load:        load,
		Spinner:     s,
		Construct:   construct,
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
	case itemsLoadedMsg:
		const defaultWidth = 20
		m.List = list.New(msg.items, ItemDelegate{}, defaultWidth, 14)
		m.List.Title = "Select a " + m.Construct
		m.List.SetShowStatusBar(false)
		m.List.SetFilteringEnabled(true)
		m.List.Styles.Title = views.TitleStyle
		m.List.Styles.PaginationStyle = views.PaginationStyle
		m.List.Styles.HelpStyle = views.HelpStyle
		m.state = stateReady
		return m, nil
	case itemsLoadFailedMsg:
		m.Err = msg.err
		m.state = stateError
		return m, tea.Quit
	case tea.WindowSizeMsg:
		if m.state != stateReady {
			return m, nil
		}
		m.List.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		if m.state != stateReady {
			return m, nil
		}
		switch keypress := msg.String(); keypress {
		case "q", "esc", "ctrl+c":
			m.Aborted = true
			m.done = true
			return m, tea.Quit
		case "enter":
			if m.List.FilterState() != list.Filtering {
				item, ok := m.List.SelectedItem().(Item)
				if ok {
					m.Choice = string(item)
					m.done = true
					return m, tea.Quit
				}
			}
		}
	}

	if m.state != stateReady {
		return m, nil
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
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
		if m.done || m.Choice != "" {
			return ""
		}
		return "\n" + m.List.View()
	}
}

func (m Model) loadCmd() tea.Cmd {
	return func() tea.Msg {
		items, err := m.Load()
		if err != nil {
			return itemsLoadFailedMsg{err: err}
		}
		return itemsLoadedMsg{items: items}
	}
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)
}
