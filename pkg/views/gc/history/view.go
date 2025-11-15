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
package history

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
	{Title: "Job Name", Width: tablelist.WidthXL},
	{Title: "Job Kind", Width: tablelist.WidthM},
	{Title: "Status", Width: tablelist.WidthS},
	{Title: "Type", Width: tablelist.WidthM},
	{Title: "Next Schedule", Width: tablelist.WidthXL},
}

func GCHistory(history []*models.GCHistory) {
	var rows []table.Row

	for _, item := range history {
		t, err := utils.FormatCreatedTime(item.Schedule.NextScheduledTime.String())
		if err != nil {
			t = item.Schedule.NextScheduledTime.String()
		}

		rows = append(rows, table.Row{
			strconv.FormatInt(int64(item.ID), 10),
			item.JobName,
			item.JobKind,
			item.JobStatus,
			item.Schedule.Type,
			t,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
