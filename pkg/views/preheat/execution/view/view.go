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
	{Title: "ID", Width: tablelist.WidthS},
	{Title: "Status", Width: tablelist.WidthM},
	{Title: "Trigger", Width: tablelist.WidthM},
	{Title: "Success Rate", Width: tablelist.WidthM},
	{Title: "Start Time", Width: tablelist.WidthL},
	{Title: "End Time", Width: tablelist.WidthL},
	{Title: "Vendor", Width: tablelist.WidthM},
}

func ViewExecution(exec *models.Execution) {
	var rows []table.Row

	startTime, _ := utils.FormatCreatedTime(exec.StartTime)
	endTime := "-"
	if exec.Status != "Running" {
		endTime, _ = utils.FormatCreatedTime(exec.EndTime)
	}

	successRate := "-"
	if m := exec.Metrics; m != nil && m.TaskCount > 0 {
		successRate = fmt.Sprintf("%d%%", m.SuccessTaskCount*100/m.TaskCount)
	}

	rows = append(rows, table.Row{
		strconv.FormatInt(exec.ID, 10),
		exec.Status,
		exec.Trigger,
		successRate,
		startTime,
		endTime,
		exec.VendorType,
	})

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
