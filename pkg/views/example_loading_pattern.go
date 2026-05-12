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

// Package views contains the reference loading state pattern for Issue #821.
// This file (example_loading_pattern.go) is a TEMPORARY template used during
// the Bubbletea v2 refactor. It will be DELETED once all model migrations
// (Steps 5–8) are complete.
//
// The pattern demonstrated here:
//  1. Init() fires both the API fetch and the spinner tick simultaneously.
//  2. Update() handles the dataLoadedMsg to clear the loading flag.
//  3. View() branches on m.loading / m.err / data-ready.
//
// Copy this pattern into each target model, replacing ExampleRow with the
// model-specific data type and someAPICall() with the real API function.
package views

import (
	"fmt"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ---------------------------------------------------------------------------
// Data types
// ---------------------------------------------------------------------------

// ExampleRow represents one row of data loaded from the API.
// Replace with the actual row type in each target model.
type ExampleRow struct {
	ID   string
	Name string
}

// exampleDataLoadedMsg is sent by fetchExampleCmd when the API call finishes.
// The err field carries any error so the TUI can show it instead of panicking.
type exampleDataLoadedMsg struct {
	data []ExampleRow
	err  error
}

// ---------------------------------------------------------------------------
// Commands
// ---------------------------------------------------------------------------

// fetchExampleCmd performs the async API call and returns an exampleDataLoadedMsg.
// Replace someAPICall() with the actual Harbor API call in each target model.
func fetchExampleCmd() tea.Cmd {
	return func() tea.Msg {
		// Simulate an API call — replace with real Harbor client call, e.g.:
		//   resp, err := api.ListProjects(...)
		//   return exampleDataLoadedMsg{data: toExampleRows(resp), err: err}
		data := []ExampleRow{
			{ID: "1", Name: "library"},
			{ID: "2", Name: "proxy-cache"},
		}
		return exampleDataLoadedMsg{data: data, err: nil}
	}
}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

// ExampleModel is the reference Bubbletea v2 model with a loading state.
// It satisfies the tea.Model interface.
type ExampleModel struct {
	// loading is true while the API call is in flight.
	loading bool

	// err holds any error returned by the API call.
	err error

	// data holds the rows fetched from the API.
	data []ExampleRow

	// spinner animates while loading is true.
	spinner spinner.Model
}

// NewExampleModel constructs an ExampleModel ready to be passed to tea.NewProgram.
func NewExampleModel() ExampleModel {
	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
	)
	return ExampleModel{
		loading: true,
		spinner: s,
	}
}

// ---------------------------------------------------------------------------
// tea.Model interface
// ---------------------------------------------------------------------------

// Init fires two commands concurrently:
//   - fetchExampleCmd: starts the async API call.
//   - a wrapper around m.spinner.Tick(): starts the spinner animation.
//
// In bubbles v2, spinner.Tick() returns tea.Msg directly. To use it with
// tea.Batch, wrap it in a func() tea.Msg — that is the definition of tea.Cmd.
func (m ExampleModel) Init() tea.Cmd {
	return tea.Batch(
		fetchExampleCmd(),
		func() tea.Msg { return m.spinner.Tick() },
	)
}

// Update handles all incoming messages.
func (m ExampleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// -----------------------------------------------------------------------
	// Data loaded from the API — clear the loading flag and store results.
	// -----------------------------------------------------------------------
	case exampleDataLoadedMsg:
		m.loading = false
		m.data = msg.data
		m.err = msg.err
		return m, nil

	// -----------------------------------------------------------------------
	// Spinner tick — forward to the spinner sub-model so it can advance its
	// frame. Only relevant while loading == true; ignored afterwards.
	// -----------------------------------------------------------------------
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	// -----------------------------------------------------------------------
	// Keyboard input — only active once data has loaded.
	// -----------------------------------------------------------------------
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the current state of the model.
//
// Three branches:
//  1. Loading  → show spinner + "Loading…" text.
//  2. Error    → show a formatted error message.
//  3. Ready    → render the actual table/list content.
func (m ExampleModel) View() tea.View {
	if m.loading {
		return tea.NewView(fmt.Sprintf("\n  %s Loading…\n", m.spinner.View()))
	}

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		return tea.NewView(fmt.Sprintf("\n  %s\n\n  Press q to quit.\n",
			errStyle.Render("Error: "+m.err.Error())))
	}

	// Data is ready — render the table.
	return tea.NewView(renderExampleTable(m.data))
}

// ---------------------------------------------------------------------------
// Rendering helpers
// ---------------------------------------------------------------------------

// renderExampleTable formats the loaded rows as a simple text table.
// In real models this is replaced by the tablelist.NewModel / tablegrid view.
func renderExampleTable(rows []ExampleRow) string {
	headerStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	out := fmt.Sprintf("\n  %s\n\n", headerStyle.Render("Example Data"))
	for _, r := range rows {
		out += fmt.Sprintf("  %-4s  %s\n", r.ID, r.Name)
	}
	out += "\n  Press q to quit.\n"
	return out
}
