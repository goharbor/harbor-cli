package get

import (
    "fmt"
    "os"

    "github.com/charmbracelet/bubbles/table"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/goharbor/go-client/pkg/sdk/v2.0/models"
    "github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
    {Title: "Field", Width: 20},
    {Title: "Value", Width: 40},
}

func DisplayUserGroup(group *models.UserGroup) {
    rows := []table.Row{
        {"ID", fmt.Sprintf("%d", group.ID)},
        {"Group Name", group.GroupName},
        {"Group Type", getGroupTypeString(group.GroupType)},
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