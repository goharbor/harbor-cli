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
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("43"))
	itemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	blockStyle    = lipgloss.NewStyle().
			Background(lipgloss.Color("81")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 1, 0)

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type Model struct {
	ready    bool
	viewport viewport.Model
	choices  []models.Permission
	cursor   int
	selected map[int]struct{}
	selects  *[]models.Permission

	loading  bool
	err      error
	spinner  spinner.Model
	fetchCmd tea.Cmd
}

type DataLoadedMsg struct {
	Choices []models.Permission
	Err     error
}

type FetchFn func() ([]models.Permission, error)

func FetchCmd(fn FetchFn) tea.Cmd {
	return func() tea.Msg {
		choices, err := fn()
		return DataLoadedMsg{Choices: choices, Err: err}
	}
}

func (m Model) Init() tea.Cmd {
	if m.loading {
		return tea.Batch(m.fetchCmd, m.spinner.Tick)
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case DataLoadedMsg:
		m.loading = false
		m.err = msg.Err
		if msg.Err == nil {
			m.choices = msg.Choices
		}
		return m, nil
	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "y":
			m.GetSelectedPermissions()
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", "space":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(msg.Height-verticalMarginHeight))
			m.viewport.SetYOffset(headerHeight)
			m.viewport.SetContent(m.listView())
			m.ready = true
			m.viewport.SetYOffset(headerHeight - 1)
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - verticalMarginHeight - 1)
		}
	}

	m.viewport.SetContent(m.listView())
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	if m.loading {
		return tea.NewView(m.spinner.View() + " Loading permissions...\n")
	}
	if m.err != nil {
		return tea.NewView("Error: " + m.err.Error() + "\n")
	}
	if !m.ready {
		return tea.NewView("\n  Initializing...")
	}
	return tea.NewView(fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView()))
}

func (m Model) headerView() string {
	title := titleStyle.Render("Select Permissions for Robot Account")
	line := strings.Repeat("─", max(0, m.viewport.Width()-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(
		fmt.Sprint(
			"  up/down: navigate • ", "enter: select permissions • ", "q: quit • ", " y: confirm\t",
		),
	)
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width()-lipgloss.Width(info)-lipgloss.Width(help)))
	return lipgloss.JoinHorizontal(lipgloss.Center, help, line, info)
}

func (m Model) listView() string {
	s := "Select Robot Permissions\n\n"
	var prev string
	for i, choice := range m.choices {
		// Render the row ith appropriate action message
		choiceRes := choice.Resource
		choiceAct := choice.Action
		now := choice.Resource
		if prev != now {
			prev = now
			s += blockStyle.Render(prev)
			s += "\n\n"
		}
		cursor := " " // no cursor
		if m.cursor == i {
			choiceRes = itemStyle.Render(choice.Resource)
			choiceAct = itemStyle.Render(choice.Action)
			cursor = ">" // cursor!
		}
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			choiceRes = selectedStyle.Render(choice.Resource)
			choiceAct = selectedStyle.Render(choice.Action)
			checked = "x" // selected!
		}
		s += fmt.Sprintf(
			"%s [%s] %s %s\n\n",
			cursor,
			checked,
			choiceAct,
			choiceRes,
		)
	}
	s += "\nPress q to quit.\n"

	return s
}

func (m Model) GetSelectedPermissions() *[]models.Permission {
	selectedPermissions := make([]models.Permission, 0, len(m.selected))
	for index := range m.selected {
		selectedPermissions = append(selectedPermissions, m.choices[index])
	}
	*m.selects = selectedPermissions
	return m.selects
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func NewModel(choices []models.Permission, selects *[]models.Permission) Model {
	return Model{
		choices:  choices,
		selected: make(map[int]struct{}),
		selects:  selects,
		loading:  false,
	}
}

func NewModelWithFetch(fetchFn FetchFn, selects *[]models.Permission) Model {
	m := NewModel(nil, selects)
	m.loading = true
	m.fetchCmd = FetchCmd(fetchFn)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m.spinner = s

	return m
}
