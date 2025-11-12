package get

import (
    "fmt"
    "os"

    "github.com/charmbracelet/bubbles/table"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/goharbor/go-client/pkg/sdk/v2.0/models"
    "github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

func DisplayUserGroup(group *models.UserGroup) {
    columns := []table.Column{
        {Title: "Group Name", Width: 30},
        {Title: "Group Type", Width: 20},
    }

    rows := []table.Row{
        {group.GroupName, getGroupTypeString(group.GroupType)},
    }

    m := tablelist.NewModel(columns, rows, len(rows))
    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }
}

func getGroupTypeString(groupType int64) string {
    switch groupType {
    case 1:
        return "LDAP"
    case 2:
        return "HTTP"
    case 3:
        return "OIDC"
    default:
        return "Unknown"
    }
}