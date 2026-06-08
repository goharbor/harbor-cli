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
	{Title: "Task ID", Width: tablelist.WidthS},
	{Title: "Exec ID", Width: tablelist.WidthS},
	{Title: "Status", Width: tablelist.WidthS},
	{Title: "Artifact", Width: tablelist.WidthXL},
	{Title: "Digest", Width: tablelist.WidthXL},
	{Title: "Type", Width: tablelist.WidthS},
	{Title: "Start Time", Width: tablelist.WidthL},
	{Title: "End Time", Width: tablelist.WidthL},
}

// Go's zero time formatted as RFC3339, used by Harbor for unset time fields.
const zeroTime = "0001-01-01T00:00:00Z"

func ListTasks(tasks []*models.Task) {
	var rows []table.Row
	for _, task := range tasks {
		startTime := "-"
		if task.StartTime != zeroTime {
			startTime, _ = utils.FormatCreatedTime(task.StartTime)
		}
		endTime := "-"
		if task.EndTime != zeroTime {
			endTime, _ = utils.FormatCreatedTime(task.EndTime)
		}

		artifact := getExtraAttr(task.ExtraAttrs, "artifact")
		digest := getExtraAttr(task.ExtraAttrs, "digest")
		kind := getExtraAttr(task.ExtraAttrs, "kind")

		rows = append(rows, table.Row{
			strconv.FormatInt(task.ID, 10),
			strconv.FormatInt(task.ExecutionID, 10),
			task.Status,
			artifact,
			digest,
			kind,
			startTime,
			endTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getExtraAttr(attrs map[string]interface{}, key string) string {
	if attrs == nil {
		return "-"
	}
	if v, ok := attrs[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return "-"
}
