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
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

// table to represent Metrics struct
var columns = []table.Column{
	{Title: "Running", Width: tablelist.WidthM},
	{Title: "Success", Width: tablelist.WidthM},
	{Title: "Error", Width: tablelist.WidthM},
	{Title: "Completed", Width: tablelist.WidthM},
	{Title: "Total", Width: tablelist.WidthM},
	{Title: "Ongoing", Width: tablelist.WidthM},
	{Title: "Trigger", Width: tablelist.WidthM},
}

func ViewScanMetrics(metrics *models.Stats) {
	var rows []table.Row
	rows = append(rows, table.Row{
		strconv.FormatInt(int64(metrics.Metrics["Running"]), 10),
		strconv.FormatInt(int64(metrics.Metrics["Success"]), 10),
		strconv.FormatInt(int64(metrics.Metrics["Error"]), 10),
		strconv.FormatInt(metrics.Completed, 10),
		strconv.FormatInt(metrics.Total, 10),
		strconv.FormatBool(metrics.Ongoing),
		metrics.Trigger,
	})

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
