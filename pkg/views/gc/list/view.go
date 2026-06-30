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
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: tablelist.WidthXS},
	{Title: "Job Name", Width: tablelist.WidthL},
	{Title: "Status", Width: tablelist.WidthM},
	{Title: "Kind", Width: tablelist.WidthM},
	{Title: "Parameters", Width: tablelist.WidthXL},
	{Title: "Creation Time", Width: tablelist.WidthXL},
}

func ListGCHistory(history []*models.GCHistory) {
	var rows []table.Row
	for _, run := range history {
		createdTime, _ := utils.FormatCreatedTime(run.CreationTime.String())
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", run.ID),
			run.JobName,
			run.JobStatus,
			run.JobKind,
			formatParams(run.JobParameters),
			createdTime,
		})
	}

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
	sort.Strings(parts)
	return strings.Join(parts, ", ")
}
