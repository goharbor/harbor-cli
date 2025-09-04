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
	{Title: "ID", Width: 4},
	{Title: "Member Name", Width: 12},
	{Title: "Type", Width: 8},
	{Title: "Role Name", Width: 16},
	{Title: "Role ID", Width: 8},
	{Title: "Project ID", Width: 12},
}

func ListMembers(members []*models.ProjectMemberEntity, wide bool) {
	var rows []table.Row
	for _, member := range members {
		memberID := strconv.FormatInt(member.ID, 10)
		roleName := utils.CamelCaseToHR(member.RoleName)

		memberType := member.EntityType
		if memberType == "u" {
			memberType = "User"
		} else if memberType == "g" {
			memberType = "Group"
		}

		if wide {
			roleID := strconv.FormatInt(member.RoleID, 10)
			projectID := strconv.FormatInt(member.ProjectID, 10)

			rows = append(rows, table.Row{
				memberID,
				member.EntityName,
				memberType,
				roleName,
				roleID,
				projectID,
			})
		} else {
			colsToRemove := []string{"Role ID", "Project ID"}
			columns = utils.RemoveColumns(columns, colsToRemove)
			rows = append(rows, table.Row{
				memberID, // Member Name
				member.EntityName,
				memberType,
				roleName,
			})
		}
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
