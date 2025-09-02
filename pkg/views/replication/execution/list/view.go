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
package list

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
	{Title: "Exec ID", Width: tablelist.WidthS},
	{Title: "Pol ID", Width: tablelist.WidthS},
	{Title: "Status", Width: tablelist.WidthL},
	{Title: "Succeed", Width: tablelist.WidthM},
	{Title: "Stopped", Width: tablelist.WidthM},
	{Title: "Total", Width: tablelist.WidthM},
	{Title: "Trigger", Width: tablelist.WidthL},
	{Title: "Start Time", Width: tablelist.WidthL},
	{Title: "End Time", Width: tablelist.WidthL},
}

func ListExecutions(executions []*models.ReplicationExecution) {
	var rows []table.Row
	for _, exec := range executions {
		createdTime, _ := utils.FormatCreatedTime(exec.StartTime.String())
		var endTime string = "-"
		if exec.Status != "InProgress" {
			endTime, _ = utils.FormatCreatedTime(exec.EndTime.String())
		}
		rows = append(rows, table.Row{
			strconv.FormatInt(exec.ID, 10),
			strconv.FormatInt(exec.PolicyID, 10),
			exec.Status,
			strconv.FormatInt(exec.Succeed, 10),
			strconv.FormatInt(exec.Stopped, 10),
			strconv.FormatInt(exec.Total, 10),
			exec.Trigger,
			createdTime,
			endTime,
		})
	}
	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
