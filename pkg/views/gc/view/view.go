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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Type", Width: tablelist.WidthM},
	{Title: "Cron", Width: tablelist.WidthL},
	{Title: "Next Scheduled Time", Width: tablelist.WidthXXL},
	{Title: "Parameters", Width: tablelist.WidthXL},
}

func ViewGCSchedule(gcHistory *models.GCHistory) {
	var rows []table.Row
	var scheduleType, scheduleCron, nextScheduledTime string
	if gcHistory != nil && gcHistory.Schedule != nil {
		scheduleType = gcHistory.Schedule.Type
		scheduleCron = gcHistory.Schedule.Cron
		nextScheduledTime = gcHistory.Schedule.NextScheduledTime.String()
	} else {
		scheduleType = "None"
		scheduleCron = "-"
		nextScheduledTime = "-"
	}

	paramsStr := ""
	if gcHistory != nil {
		paramsStr = gcHistory.JobParameters
	}

	rows = append(rows, table.Row{
		scheduleType,
		scheduleCron,
		nextScheduledTime,
		formatParams(paramsStr),
	})

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func formatParams(paramsStr string) string {
	if paramsStr == "" {
		return "-"
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(paramsStr), &m); err != nil {
		return paramsStr
	}
	var parts []string
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, ", ")
}
