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

func memberColumns(wide bool) []table.Column {
	if wide {
		return columns
	}
	colsToRemove := []string{"Role ID", "Project ID"}
	return utils.RemoveColumns(columns, colsToRemove)
}

func memberRow(member *models.ProjectMemberEntity, wide bool) table.Row {
	memberID := strconv.FormatInt(member.ID, 10)
	roleName := utils.CamelCaseToHR(member.RoleName)

	memberType := member.EntityType
	switch memberType {
	case "u":
		memberType = "User"
	case "g":
		memberType = "Group"
	}

	if wide {
		roleID := strconv.FormatInt(member.RoleID, 10)
		projectID := strconv.FormatInt(member.ProjectID, 10)

		return table.Row{
			memberID,
			member.EntityName,
			memberType,
			roleName,
			roleID,
			projectID,
		}
	}

	return table.Row{
		memberID,
		member.EntityName,
		memberType,
		roleName,
	}
}

func ViewMember(member *models.ProjectMemberEntity, wide bool) {
	rows := []table.Row{memberRow(member, wide)}

	m := tablelist.NewModel(memberColumns(wide), rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
