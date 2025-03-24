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

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "ID", Width: 5},
	{Title: "Name", Width: 10},
	{Title: "Status", Width: 10},
	{Title: "Endpoint URL", Width: 20},
	{Title: "Provider", Width: 15},
	{Title: "Creation Time", Width: 15},
	{Title: "Description", Width: 20},
}

func ViewRegistry(registry *models.Registry) {
	var rows []table.Row
	createdTime, _ := utils.FormatCreatedTime(registry.CreationTime.String())
	rows = append(rows, table.Row{
		fmt.Sprintf("%d", registry.ID),
		registry.Name,
		registry.Status,
		registry.URL,
		registry.Type,
		createdTime,
		registry.Description,
	})

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
