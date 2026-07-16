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
package interactive

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/spf13/cobra"
)

const (
	defaultListWidth  = 40
	defaultListHeight = 12
	minListWidth      = 30
	minListHeight     = 6
	verticalChrome    = 8
)

type commandItem struct {
	title       string
	description string
	command     *cobra.Command
}

func (i commandItem) Title() string       { return i.title }
func (i commandItem) Description() string { return i.description }
func (i commandItem) FilterValue() string { return i.title + " " + i.description }

type stackEntry struct {
	command *cobra.Command
	index   int
}

type Model struct {
	list     list.Model
	stack    []stackEntry
	current  *cobra.Command
	selected *cobra.Command
	width    int
	height   int
}

type ItemDelegate struct{}

func (d ItemDelegate) Height() int  { return 2 }
func (d ItemDelegate) Spacing() int { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(commandItem)
	if !ok {
		return
	}

	title := item.title
	description := item.description
	if description == "" {
		description = "No description available"
	}

	titleStyle := views.ItemStyle
	descriptionStyle := views.GrayStyle.PaddingLeft(4)
	if index == m.Index() {
		titleStyle = views.SelectedItemStyle
		descriptionStyle = views.GrayStyle.PaddingLeft(2)
		title = "> " + title
	}

	fmt.Fprint(w, titleStyle.Render(title))
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, descriptionStyle.Render(description))
}

func NewModel(root *cobra.Command) Model {
	items := commandItems(root)
	l := list.New(items, ItemDelegate{}, defaultListWidth, defaultListHeight)
	l.Title = fmt.Sprintf("Interactive Harbor: %s", root.CommandPath())
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = views.TitleStyle
	l.Styles.PaginationStyle = views.PaginationStyle
	l.Styles.HelpStyle = views.HelpStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return nil
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return nil
	}

	return Model{
		list:    l,
		current: root,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(dynamicListSize(msg.Width, msg.Height))
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.list.FilterState() == list.Filtering {
				break
			}
			selectedItem, ok := m.list.SelectedItem().(commandItem)
			if !ok {
				return m, nil
			}
			if len(commandItems(selectedItem.command)) == 0 {
				m.selected = selectedItem.command
				return m, tea.Quit
			}

			m.stack = append(m.stack, stackEntry{command: m.current, index: m.list.Index()})
			m.current = selectedItem.command
			m.list.SetItems(commandItems(selectedItem.command))
			m.list.Title = fmt.Sprintf("Interactive Harbor: %s", selectedItem.command.CommandPath())
			m.list.ResetSelected()
			m.list.ResetFilter()
			return m, nil
		case "esc":
			if m.list.FilterState() == list.Filtering {
				break
			}
			if len(m.stack) == 0 {
				return m, tea.Quit
			}
			last := m.stack[len(m.stack)-1]
			m.stack = m.stack[:len(m.stack)-1]
			m.current = last.command
			m.list.SetItems(commandItems(last.command))
			m.list.Title = fmt.Sprintf("Interactive Harbor: %s", last.command.CommandPath())
			m.list.Select(last.index)
			m.list.ResetFilter()
			return m, nil
		case "q", "ctrl+c":
			if m.list.FilterState() != list.Filtering {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.selected != nil {
		return renderSelectionSummary(m.selected)
	}

	breadcrumb := views.BlueStyle.Render("Path: " + strings.Join(commandPathSlice(m.current), " > "))
	hint := views.GrayStyle.Render("Enter: open/select  Esc: back  /: filter  q: quit")
	body := "\n" + m.list.View() + "\n" + breadcrumb + "\n" + hint + "\n"
	if m.width == 0 {
		return body
	}

	return lipgloss.NewStyle().Width(m.width).Render(body)
}

func Run(root *cobra.Command) (*cobra.Command, error) {
	program := tea.NewProgram(NewModel(root), tea.WithAltScreen())
	finalModel, err := program.Run()
	if err != nil {
		return nil, err
	}

	model, ok := finalModel.(Model)
	if !ok {
		return nil, fmt.Errorf("unexpected interactive result")
	}

	return model.selected, nil
}

func renderSelectionSummary(cmd *cobra.Command) string {
	var sections []string
	sections = append(sections, views.BoldStyle.Render("Selected command"))
	sections = append(sections, views.GreenStyle.Render(cmd.CommandPath()))

	if cmd.Short != "" {
		sections = append(sections, "")
		sections = append(sections, views.BoldStyle.Render("Summary"))
		sections = append(sections, cmd.Short)
	}

	if cmd.Long != "" && cmd.Long != cmd.Short {
		sections = append(sections, "")
		sections = append(sections, views.BoldStyle.Render("Details"))
		sections = append(sections, strings.TrimSpace(cmd.Long))
	}

	if cmd.Example != "" {
		sections = append(sections, "")
		sections = append(sections, views.BoldStyle.Render("Examples"))
		sections = append(sections, strings.TrimSpace(cmd.Example))
	}

	if cmd.HasAvailableFlags() {
		sections = append(sections, "")
		sections = append(sections, views.BoldStyle.Render("Flags"))
		sections = append(sections, "This command defines flags. Run the command with --help to inspect full flag details.")
	}

	sections = append(sections, "")
	sections = append(sections, views.YellowStyle.Render("v1 preview: interactive mode currently browses the command tree and shows the selected leaf command."))

	return "\n" + strings.Join(sections, "\n") + "\n"
}

func commandItems(cmd *cobra.Command) []list.Item {
	commands := visibleSubcommands(cmd)
	items := make([]list.Item, 0, len(commands))
	for _, sub := range commands {
		items = append(items, commandItem{
			title:       sub.Name(),
			description: sub.Short,
			command:     sub,
		})
	}
	return items
}

func visibleSubcommands(cmd *cobra.Command) []*cobra.Command {
	commands := cmd.Commands()
	items := make([]*cobra.Command, 0, len(commands))
	for _, sub := range commands {
		if !sub.IsAvailableCommand() || sub.Hidden {
			continue
		}
		if sub.Name() == "help" || sub.Name() == "completion" || sub.Name() == "interactive" {
			continue
		}
		items = append(items, sub)
	}
	return items
}

func commandPathSlice(cmd *cobra.Command) []string {
	if cmd == nil {
		return nil
	}
	return strings.Fields(cmd.CommandPath())
}

func dynamicListSize(width, height int) (int, int) {
	listWidth := width
	if listWidth < minListWidth {
		listWidth = minListWidth
	}

	listHeight := height - verticalChrome
	if listHeight < minListHeight {
		listHeight = minListHeight
	}

	return listWidth, listHeight
}
