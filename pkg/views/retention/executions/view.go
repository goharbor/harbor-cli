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
package executions

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"golang.org/x/term"
)

// ErrUserAborted is returned when the user quits the selection without
// picking an execution.
var ErrUserAborted = errors.New("user aborted selection")

func truncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 160
	}
	return width
}

func getAdjustedColumns() []table.Column {
	totalWidth := getTerminalWidth()
	columnWidths := []int{totalWidth / 12, totalWidth / 12, totalWidth / 12, totalWidth / 10, totalWidth / 8, totalWidth / 8, totalWidth / 8}
	return []table.Column{
		{Title: "ID", Width: columnWidths[0]},
		{Title: "Policy", Width: columnWidths[1]},
		{Title: "Dry Run", Width: columnWidths[2]},
		{Title: "Trigger", Width: columnWidths[3]},
		{Title: "Status", Width: columnWidths[4]},
		{Title: "Start Time", Width: columnWidths[5]},
		{Title: "End Time", Width: columnWidths[6]},
	}
}

func ListRetentionExecutions(executions []*models.RetentionExecution) {
	var rows []table.Row
	columns := getAdjustedColumns()

	for _, execution := range executions {
		rows = append(rows, table.Row{
			truncateString(fmt.Sprintf("%d", execution.ID), columns[0].Width),
			truncateString(fmt.Sprintf("%d", execution.PolicyID), columns[1].Width),
			truncateString(fmt.Sprintf("%v", execution.DryRun), columns[2].Width),
			truncateString(execution.Trigger, columns[3].Width),
			truncateString(execution.Status, columns[4].Width),
			truncateString(execution.StartTime, columns[5].Width),
			truncateString(execution.EndTime, columns[6].Width),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// RetentionExecutionList prompts the user to select one of the executions.
// It returns ErrUserAborted when the user quits without choosing.
func RetentionExecutionList(executions []*models.RetentionExecution) (int64, error) {
	itemsList := make([]list.Item, len(executions))
	items := map[string]int64{}

	for i, execution := range executions {
		display := fmt.Sprintf("ID: %d | Policy: %d | Status: %s | Trigger: %s | Dry Run: %v | Start: %s | End: %s",
			execution.ID,
			execution.PolicyID,
			execution.Status,
			execution.Trigger,
			execution.DryRun,
			execution.StartTime,
			execution.EndTime,
		)
		items[display] = execution.ID
		itemsList[i] = selection.Item(display)
	}

	m := selection.NewModel(itemsList, "a Retention Execution")
	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running selection UI:", err)
		os.Exit(1)
	}

	if model, ok := p.(selection.Model); ok {
		if model.Aborted {
			return 0, ErrUserAborted
		}
		id, ok := items[model.Choice]
		if !ok {
			return 0, errors.New("no execution selected")
		}
		return id, nil
	}
	return 0, errors.New("unexpected program result")
}
