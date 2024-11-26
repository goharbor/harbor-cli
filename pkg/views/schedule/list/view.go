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
	{Title: "Cron", Width: 18},
	{Title: "Vendor Type", Width: 28},
	{Title: "Update Time", Width: 20},
}

func ListSchedule(schedule []*models.ScheduleTask) {
	var rows []table.Row
	for _, regis := range schedule {
		updatedTime, _ := utils.FormatCreatedTime(regis.UpdateTime.String())
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", regis.ID),
			regis.Cron,
			regis.VendorType,
			updatedTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
