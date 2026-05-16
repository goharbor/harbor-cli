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

package logviewer

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogMsg string
type triggerFetchMsg struct{}

type Model struct {
	viewport    viewport.Model
	jobID       string
	fetcher     func(string) (string, error)
	follow      bool
	interval    time.Duration
	lastContent string
	ready       bool
	isFetching  bool
	err         error
}

func NewModel(jobID string, fetcher func(string) (string, error), follow bool, interval time.Duration) Model {
	return Model{
		jobID:    jobID,
		fetcher:  fetcher,
		follow:   follow,
		interval: interval,
	}
}

func (m Model) Init() tea.Cmd {
	return m.fetchLogCmd
}

func (m Model) fetchLogCmd() tea.Msg {
	content, err := m.fetcher(m.jobID)
	if err != nil {
		return err
	}
	return LogMsg(content)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := 2
		footerHeight := 2
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	case LogMsg:
		m.isFetching = false
		content := string(msg)
		if content != m.lastContent {
			m.viewport.SetContent(content)
			m.lastContent = content
			m.viewport.GotoBottom()
		}
		if m.follow {
			return m, tea.Tick(m.interval, func(t time.Time) tea.Msg {
				return triggerFetchMsg{}
			})
		}

	case triggerFetchMsg:
		if !m.isFetching {
			m.isFetching = true
			return m, m.fetchLogCmd
		}

	case error:
		m.err = msg
		m.isFetching = false
		return m, tea.Quit
	}

	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	if m.err != nil {
		return fmt.Sprintf("\n  Error: %v\n", m.err)
	}

	header := lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf(" Logs for Job: %s", m.jobID))
	footer := lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(" %.0f%% • q to quit", m.viewport.ScrollPercent()*100))

	return fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), footer)
}
