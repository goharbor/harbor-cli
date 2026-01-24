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

package gc

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

type GCJobParams struct {
	DryRun         bool `json:"dry_run"`
	DeleteUntagged bool `json:"delete_untagged"`
}

var columns = []table.Column{
	{Title: "ID", Width: 10},
	{Title: "Status", Width: 15},
	{Title: "Dry Run", Width: 10},
	{Title: "Creation Time", Width: 25},
	{Title: "Update Time", Width: 25},
}

func ListGC(history []*models.GCHistory) {
	var rows []table.Row
	for _, job := range history {
		creationTime, _ := utils.FormatCreatedTime(job.CreationTime.String())
		updateTime, _ := utils.FormatCreatedTime(job.UpdateTime.String())
		dryRun := "false"

		if job.JobParameters != "" {
			var params GCJobParams
			if err := json.Unmarshal([]byte(job.JobParameters), &params); err == nil {
				dryRun = strconv.FormatBool(params.DryRun)
			}
		}

		// Note: JobParameters is usually a JSON string. For simplicity we display it as is or handle parsing if needed.
		// Usually contains {"dry_run": true/false}

		rows = append(rows, table.Row{
			strconv.FormatInt(job.ID, 10),
			job.JobStatus,
			dryRun,
			creationTime,
			updateTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
