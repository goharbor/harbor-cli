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
package view

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Property", Width: tablelist.WidthL},
	{Title: "Value", Width: tablelist.Width3XL},
}

var order = []string{
	"ID",
	"Status",
	"Resource Type",
	"Source Resource",
	"Destination Resource",
	"Operation",
	"Job ID",
	"Execution ID",
	"Start Time",
	"End Time",
}

func ViewTask(task *models.ReplicationTask) {
	startTime, _ := utils.FormatCreatedTime(task.StartTime.String())
	var endTime string = "-"
	if task.Status != "InProgress" {
		endTime, _ = utils.FormatCreatedTime(task.EndTime.String())
	}

	taskMap := map[string]string{
		"ID":                   strconv.FormatInt(task.ID, 10),
		"Status":               task.Status,
		"Resource Type":        task.ResourceType,
		"Source Resource":      task.SrcResource,
		"Destination Resource": task.DstResource,
		"Operation":            task.Operation,
		"Job ID":               task.JobID,
		"Execution ID":         strconv.FormatInt(task.ExecutionID, 10),
		"Start Time":           startTime,
		"End Time":             endTime,
	}

	var rows []table.Row
	for _, key := range order {
		rows = append(rows, table.Row{
			key,
			taskMap[key],
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func ListTasks(tasks []*models.ReplicationTask) {
	if len(tasks) == 0 {
		fmt.Println("No replication tasks found")
		return
	}

	if len(tasks) == 1 {
		ViewTask(tasks[0])
		return
	}

	var listColumns = []table.Column{
		{Title: "ID", Width: tablelist.WidthS},
		{Title: "Status", Width: tablelist.WidthM},
		{Title: "Resource Type", Width: tablelist.WidthM},
		{Title: "Source Resource", Width: tablelist.WidthXL},
		{Title: "Destination Resource", Width: tablelist.WidthXL},
		{Title: "Operation", Width: tablelist.WidthM},
		{Title: "Start Time", Width: tablelist.WidthL},
	}

	var rows []table.Row
	for _, task := range tasks {
		startTime, _ := utils.FormatCreatedTime(task.StartTime.String())

		rows = append(rows, table.Row{
			strconv.FormatInt(task.ID, 10),
			task.Status,
			task.ResourceType,
			task.SrcResource,
			task.DstResource,
			task.Operation,
			startTime,
		})
	}

	m := tablelist.NewModel(listColumns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
