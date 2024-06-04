package multiselect

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	content  string
	ready    bool
	viewport viewport.Model
	choices  []models.Permission // items on the to-do list
	cursor   int                 // which to-do list item our cursor is pointing at
	selected map[int]struct{}    // which to-do items are selected
	selects  *[]models.Permission
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
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
		case "enter", " ":
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
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.listView())
			m.ready = true
			m.viewport.YPosition = headerHeight - 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight - 1
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	m.viewport.SetContent(m.listView())
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
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
	}
}
