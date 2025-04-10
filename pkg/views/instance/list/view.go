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
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: 6},
	{Title: "Name", Width: 12},
	{Title: "Provider", Width: 12},
	{Title: "Endpoint", Width: 24},
	{Title: "Status", Width: 12},
	{Title: "Auth Mode", Width: 10},
	{Title: "Description", Width: 12},
	{Title: "Default", Width: 8},
	{Title: "Insecure", Width: 8},
	{Title: "Enabled", Width: 8},
	{Title: "Setup Timestamp", Width: 20},
}

func ListInstance(instance []*models.Instance) {
	var rows []table.Row
	for _, regis := range instance {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", regis.ID),
			regis.Name,
			regis.Vendor,
			regis.Endpoint,
			regis.Status,
			regis.AuthMode,
			regis.Description,
			fmt.Sprintf("%t", regis.Default),
			fmt.Sprintf("%t", regis.Insecure),
			fmt.Sprintf("%t", regis.Enabled),
			time.Unix(regis.SetupTimestamp, 0).Format("2006-01-02 15:04:05"), // updated
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
