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
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: tablelist.WidthXS},
	{Title: "Name", Width: tablelist.WidthM},
	{Title: "Provider", Width: tablelist.WidthM},
	{Title: "Endpoint", Width: tablelist.WidthXXL},
	{Title: "Status", Width: tablelist.WidthM},
	{Title: "Auth Mode", Width: tablelist.WidthM},
	{Title: "Description", Width: tablelist.WidthM},
	{Title: "Default", Width: tablelist.WidthS},
	{Title: "Insecure", Width: tablelist.WidthS},
	{Title: "Enabled", Width: tablelist.WidthS},
	{Title: "Setup Timestamp", Width: tablelist.WidthXL},
}

func ViewInstance(instance *models.Instance) {
	var rows []table.Row
	rows = append(rows, table.Row{
		fmt.Sprintf("%d", instance.ID),
		instance.Name,
		instance.Vendor,
		instance.Endpoint,
		instance.Status,
		instance.AuthMode,
		instance.Description,
		fmt.Sprintf("%t", instance.Default),
		fmt.Sprintf("%t", instance.Insecure),
		fmt.Sprintf("%t", instance.Enabled),
		time.Unix(instance.SetupTimestamp, 0).Format("2006-01-02 15:04:05"),
	})

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
