package search

import (
    "fmt"
    "os"
    "strconv"

    "github.com/charmbracelet/bubbles/table"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/goharbor/go-client/pkg/sdk/v2.0/client/usergroup"
    "github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
    {Title: "ID", Width: 10},
    {Title: "Group Name", Width: 30},
    {Title: "Group Type", Width: 15},
}

func DisplayUserGroupSearchResults(resp *usergroup.SearchUserGroupsOK) {
    var rows []table.Row
    for _, group := range resp.Payload {
        groupType := getGroupTypeString(group.GroupType)
        rows = append(rows, table.Row{
            strconv.Itoa(int(group.ID)),
            group.GroupName,
            groupType,
        })
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