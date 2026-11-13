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
package multiselectv2

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views"
	"golang.org/x/term"
)

type state string

const (
	stateLoading state = "loading"
	stateReady   state = "ready"
	stateError   state = "error"
)

type choicesLoadedMsg struct {
	choices []models.Permission
}

type choicesLoadFailedMsg struct {
	err error
}

type Loader func() ([]models.Permission, error)

type Model struct {
	Load        Loader
	Selects     *[]models.Permission
	Spinner     spinner.Model
	Err         error
	LoadingText string
	windowSize  *tea.WindowSizeMsg
	viewport    viewport.Model
	choices     []models.Permission
	cursor      int
	selected    map[int]struct{}
	done        bool
	ready       bool
	animate     bool
	state       state
}

func NewModel(selects *[]models.Permission, loadingText string, load Loader) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	if loadingText == "" {
		loadingText = "Loading..."
	}

	return Model{
		Load:        load,
		Selects:     selects,
		Spinner:     s,
		LoadingText: loadingText,
		selected:    make(map[int]struct{}),
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
	case choicesLoadedMsg:
		m.choices = msg.choices
		m.state = stateReady
		if m.windowSize != nil {
			return m.handleWindowSize(*m.windowSize)
		}
		return m, nil
	case choicesLoadFailedMsg:
		m.Err = msg.err
		m.state = stateError
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.windowSize = &msg
		if m.state != stateReady {
			return m, nil
		}
		return m.handleWindowSize(msg)
	case tea.KeyMsg:
		if m.state != stateReady {
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.done = true
			return m, tea.Quit
		case "y":
			m.getSelectedPermissions()
			m.done = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.ensureCursorVisible()
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
				m.ensureCursorVisible()
			}
		case "enter", " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	if m.state != stateReady {
		return m, nil
	}

	m.viewport.SetContent(m.listView())
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
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
		if !m.ready {
			return "\n  Initializing..."
		}
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	}
}

func (m Model) loadCmd() tea.Cmd {
	return func() tea.Msg {
		choices, err := m.Load()
		if err != nil {
			return choicesLoadFailedMsg{err: err}
		}
		return choicesLoadedMsg{choices: choices}
	}
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)
}

func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	if !m.ready {
		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = false
		m.viewport.SetContent(m.listView())
		m.ready = true
		m.viewport.YPosition = headerHeight - 1
		m.ensureCursorVisible()
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight - 1
		m.ensureCursorVisible()
	}

	return m, nil
}

func (m Model) headerView() string {
	title := titleStyle.Render("Select Permissions for Robot Account")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(
		fmt.Sprint(
			"  up/down: navigate • ", "enter: select permissions • ", "q: quit • ", " y: confirm\t",
		),
	)
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)-lipgloss.Width(help)))
	return lipgloss.JoinHorizontal(lipgloss.Center, help, line, info)
}

func (m Model) listView() string {
	s := "Select Robot Permissions\n\n"
	var prev string
	for i, choice := range m.choices {
		choiceRes := choice.Resource
		choiceAct := choice.Action
		if prev != choice.Resource {
			prev = choice.Resource
			s += blockStyle.Render(prev)
			s += "\n\n"
		}

		cursor := " "
		if m.cursor == i {
			choiceRes = itemStyle.Render(choice.Resource)
			choiceAct = itemStyle.Render(choice.Action)
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			choiceRes = selectedStyle.Render(choice.Resource)
			choiceAct = selectedStyle.Render(choice.Action)
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s %s\n\n", cursor, checked, choiceAct, choiceRes)
	}
	s += "\nPress q to quit.\n"
	return s
}

func (m *Model) getSelectedPermissions() {
	selectedPermissions := make([]models.Permission, 0, len(m.selected))
	for index := range m.selected {
		selectedPermissions = append(selectedPermissions, m.choices[index])
	}
	*m.Selects = selectedPermissions
}

func (m *Model) ensureCursorVisible() {
	if !m.ready || m.viewport.Height <= 0 {
		return
	}

	cursorLine := m.cursorLineOffset()
	if cursorLine < m.viewport.YOffset {
		m.viewport.SetYOffset(cursorLine)
		return
	}

	if cursorLine >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.SetYOffset(cursorLine - m.viewport.Height + 1)
	}
}

func (m Model) cursorLineOffset() int {
	line := 2
	var prev string

	for i, choice := range m.choices {
		if prev != choice.Resource {
			prev = choice.Resource
			line += 2
		}

		if i == m.cursor {
			return line
		}

		line += 2
	}

	return max(0, line-1)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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
