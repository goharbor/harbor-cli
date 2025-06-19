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
	{Title: "ID", Width: tablelist.WidthS},
	{Title: "Name", Width: tablelist.WidthL},
	{Title: "Enabled", Width: tablelist.WidthS},
	{Title: "Source", Width: tablelist.Width3XL},
	{Title: "Destination", Width: tablelist.Width3XL},
	{Title: "Trigger Type", Width: tablelist.WidthM},
	{Title: "Override", Width: tablelist.WidthS},
	{Title: "Creation Time", Width: tablelist.WidthL},
	{Title: "Last Modified", Width: tablelist.WidthL},
}

func ListPolicies(rpolicies []*models.ReplicationPolicy) {
	var rows []table.Row
	for _, rpolicy := range rpolicies {

		createdTime, _ := utils.FormatCreatedTime(rpolicy.CreationTime.String())
		modifledTime, _ := utils.FormatCreatedTime(rpolicy.UpdateTime.String())
		rows = append(rows, table.Row{
			strconv.FormatInt(rpolicy.ID, 10),
			rpolicy.Name,
			strconv.FormatBool(rpolicy.Enabled),
			rpolicy.SrcRegistry.Name,
			rpolicy.DestRegistry.Name,
			rpolicy.Trigger.Type,
			strconv.FormatBool(rpolicy.Override),
			createdTime,
			modifledTime,
		})
	}
	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
