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
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: tablelist.WidthS},
	{Title: "Job Name", Width: tablelist.WidthXL},
	{Title: "Status", Width: tablelist.WidthM},
	{Title: "Kind", Width: tablelist.WidthM},
	{Title: "Parameters", Width: tablelist.WidthXL * 2},
	{Title: "Creation Time", Width: tablelist.WidthXL},
}

func ListGCHistory(history []*models.GCHistory) {
	var rows []table.Row
	for _, run := range history {
		createdTime, err := utils.FormatCreatedTime(run.CreationTime.String())
		if err != nil {
			fmt.Println("Error formatting created time:", err)
			os.Exit(1)
		}
		rows = append(rows, historyRows(run, createdTime)...)
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func historyRows(run *models.GCHistory, createdTime string) []table.Row {
	params := formatParams(run.JobParameters)
	rows := []table.Row{{
		fmt.Sprintf("%d", run.ID),
		run.JobName,
		run.JobStatus,
		run.JobKind,
		params[0],
		createdTime,
	}}
	for _, param := range params[1:] {
		rows = append(rows, table.Row{"", "", "", "", param, ""})
	}
	return rows
}

func formatParams(paramsStr string) []string {
	if paramsStr == "" {
		return []string{"-"}
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(paramsStr), &m); err != nil {
		return []string{paramsStr}
	}
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s: %s", k, formatParamValue(v)))
	}
	if len(parts) == 0 {
		return []string{"-"}
	}
	sort.Strings(parts)
	return parts
}

func formatParamValue(value interface{}) string {
	if value, ok := value.(string); ok {
		return value
	}

	formatted, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(formatted)
}
