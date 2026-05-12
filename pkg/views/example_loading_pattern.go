package views

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// SomeRow is a placeholder for actual data (like models.Instance, models.Project, etc.)
type SomeRow struct {
	ID   int
	Name string
}

// ExampleModel demonstrates the standard loading pattern
type ExampleModel struct {
	loading bool
	err     error
	data    []SomeRow
	spinner spinner.Model
}

// dataLoadedMsg is a generic message type for when data fetching completes
type dataLoadedMsg struct {
	data []SomeRow
	err  error
}

// simulate API call
func someAPICall() ([]SomeRow, error) {
	time.Sleep(1 * time.Second)
	return []SomeRow{{ID: 1, Name: "Row 1"}, {ID: 2, Name: "Row 2"}}, nil
}

// fetchDataCmd executes the API call asynchronously
func fetchDataCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := someAPICall()
		return dataLoadedMsg{data: data, err: err}
	}
}

// NewExampleModel initializes the model with a spinner and sets loading to true
func NewExampleModel() ExampleModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return ExampleModel{
		loading: true,
		spinner: s,
	}
}

func (m ExampleModel) Init() tea.Cmd {
	// Start both the API call and the spinner ticking
	return tea.Batch(fetchDataCmd(), m.spinner.Tick)
}

func (m ExampleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case dataLoadedMsg:
		// Data has arrived, update state
		m.loading = false
		m.data = msg.data
		m.err = msg.err
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	// Only update the spinner if we are still loading
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Update list/table/grid here when not loading
	return m, nil
}

func (m ExampleModel) View() tea.View {
	if m.loading {
		// Show loading spinner
		return tea.NewView(fmt.Sprintf("%s Loading data...\n", m.spinner.View()))
	}
	if m.err != nil {
		// Show error state
		return tea.NewView(fmt.Sprintf("Error: %v\n", m.err))
	}

	// Render the actual UI with data
	output := "Data Loaded Successfully:\n\n"
	for _, row := range m.data {
		output += fmt.Sprintf("- %d: %s\n", row.ID, row.Name)
	}
	output += "\nPress q to quit."

	return tea.NewView(output)
}
