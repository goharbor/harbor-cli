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
	{Title: "ID", Width: 6},
	{Title: "Status", Width: 10},
	{Title: "Dry Run", Width: 10},
	{Title: "Execution Type", Width: 16},
	{Title: "Start Time", Width: 18},
}

func ListRetentionRules(retention []*models.RetentionExecution) {
	var rows []table.Row
	for _, regis := range retention {
		createdTime, _ := utils.FormatCreatedTime(regis.StartTime)
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", regis.ID),
			regis.Status,
			fmt.Sprintf("%v", regis.DryRun),
			regis.Trigger,
			createdTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}