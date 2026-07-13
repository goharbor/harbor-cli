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
package tablegridv2

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablegrid"
	"golang.org/x/term"
)

type state string

const (
	stateLoading state = "loading"
	stateReady   state = "ready"
	stateError   state = "error"
)

type gridLoadedMsg struct {
	config tablegrid.Config
}

type gridLoadFailedMsg struct {
	err error
}

type Loader func() (tablegrid.Config, error)

type Model struct {
	Grid        *tablegrid.TableGrid
	Load        Loader
	Spinner     spinner.Model
	Err         error
	LoadingText string
	animate     bool
	state       state
}

func NewModel(loadingText string, load Loader) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	if loadingText == "" {
		loadingText = "Loading..."
	}

	return Model{
		Load:        load,
		Spinner:     s,
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
	case gridLoadedMsg:
		m.Grid = tablegrid.New(msg.config)
		m.state = stateReady
		return m, nil
	case gridLoadFailedMsg:
		m.Err = msg.err
		m.state = stateError
		return m, tea.Quit
	}

	if m.Grid == nil {
		return m, nil
	}

	gridModel, cmd := m.Grid.Update(msg)
	if grid, ok := gridModel.(*tablegrid.TableGrid); ok {
		m.Grid = grid
	}
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
		if m.Grid == nil {
			return views.BaseStyle.Render("")
		}
		return m.Grid.View()
	}
}

func (m Model) loadCmd() tea.Cmd {
	return func() tea.Msg {
		config, err := m.Load()
		if err != nil {
			return gridLoadFailedMsg{err: err}
		}
		return gridLoadedMsg{config: config}
	}
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)
}
