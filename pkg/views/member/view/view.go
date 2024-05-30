package view

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

type model struct {
	table table.Model
}

var columns = []table.Column{
	{Title: "ID", Width: 4},
	{Title: "Member Name", Width: 12},
	{Title: "Type", Width: 8},
	{Title: "Role Name", Width: 16},
	{Title: "Role ID", Width: 8},
	{Title: "Project ID", Width: 12},
}

func ViewMember(member *models.ProjectMemberEntity, wide bool) {
	var rows []table.Row
	memberID := strconv.FormatInt(member.ID, 10)
	projectID := strconv.FormatInt(member.ProjectID, 10)
	roleName := utils.CamelCaseToHR(member.RoleName)

	memberType := member.EntityType
	if memberType == "u" {
		memberType = "User"
	} else if memberType == "g" {
		memberType = "Group"
	}

	if wide {
		roleID := strconv.FormatInt(member.RoleID, 10)

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
		log.Println(columns)
		rows = append(rows, table.Row{
			memberID, // Member Name
			member.EntityName,
			memberType,
			roleName,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
