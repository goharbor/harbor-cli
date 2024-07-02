package list

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
	{Title: "User", Width: 18},
	{Title: "Resource", Width: 18},
	{Title: "Resource Type", Width: 14},
	{Title: "Operation", Width: 10},
	{Title: "Time", Width: 16},
}

func ListLogs(logs []*models.AuditLog) {
	var rows []table.Row
	for _, log := range logs {
		operationTime, _ := utils.FormatCreatedTime(log.OpTime.String())
		rows = append(rows, table.Row{
			log.Username,
			log.Resource,
			log.ResourceType,
			log.Operation,
			operationTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
