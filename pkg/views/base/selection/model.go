package selection

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
	List   list.Model
	Choice string
}

func NewModel(items []list.Item, construct string) Model {
	const defaultWidth = 20
	l := list.New(items, ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select a " + construct
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = views.TitleStyle
	l.Styles.PaginationStyle = views.PaginationStyle
	l.Styles.HelpStyle = views.HelpStyle

	return Model{List: l}
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
		case "enter":
			i, ok := m.List.SelectedItem().(Item)
			if ok {
				m.Choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Choice != "" {
		return ""
	}
	return "\n" + m.List.View()
}
