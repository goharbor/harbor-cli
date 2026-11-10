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
package tasks

import (
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

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 180
	}
	return width
}

func getAdjustedColumns() []table.Column {
	totalWidth := getTerminalWidth()
	columnWidths := []int{totalWidth / 15, totalWidth / 10, totalWidth / 7, totalWidth / 10, totalWidth / 10, totalWidth / 10, totalWidth / 10, totalWidth / 10, totalWidth / 10, totalWidth / 10}
	return []table.Column{
		{Title: "ID", Width: columnWidths[0]},
		{Title: "Exec ID", Width: columnWidths[1]},
		{Title: "Repository", Width: columnWidths[2]},
		{Title: "Status", Width: columnWidths[3]},
		{Title: "Retained", Width: columnWidths[4]},
		{Title: "Total", Width: columnWidths[5]},
		{Title: "Code", Width: columnWidths[6]},
		{Title: "Revision", Width: columnWidths[7]},
		{Title: "Start Time", Width: columnWidths[8]},
		{Title: "End Time", Width: columnWidths[9]},
	}
}

func ListRetentionTasks(tasks []*models.RetentionExecutionTask) {
	var rows []table.Row
	columns := getAdjustedColumns()

	for _, task := range tasks {
		rows = append(rows, table.Row{
			truncateString(fmt.Sprintf("%d", task.ID), columns[0].Width),
			truncateString(fmt.Sprintf("%d", task.ExecutionID), columns[1].Width),
			truncateString(task.Repository, columns[2].Width),
			truncateString(task.Status, columns[3].Width),
			truncateString(fmt.Sprintf("%d", task.Retained), columns[4].Width),
			truncateString(fmt.Sprintf("%d", task.Total), columns[5].Width),
			truncateString(fmt.Sprintf("%d", task.StatusCode), columns[6].Width),
			truncateString(fmt.Sprintf("%d", task.StatusRevision), columns[7].Width),
			truncateString(task.StartTime, columns[8].Width),
			truncateString(task.EndTime, columns[9].Width),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func RetentionTaskList(tasks []*models.RetentionExecutionTask, choice chan<- int64) {
	itemsList := make([]list.Item, len(tasks))
	items := map[string]int64{}

	for i, task := range tasks {
		display := fmt.Sprintf("ID: %d | Exec: %d | Repo: %s | Status: %s | Retained: %d | Total: %d | Code: %d",
			task.ID,
			task.ExecutionID,
			task.Repository,
			task.Status,
			task.Retained,
			task.Total,
			task.StatusCode,
		)
		items[display] = task.ID
		itemsList[i] = selection.Item(display)
	}

	m := selection.NewModel(itemsList, "a Retention Task")
	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running selection UI:", err)
		os.Exit(1)
	}

	if model, ok := p.(selection.Model); ok {
		choice <- items[model.Choice]
	}
}
