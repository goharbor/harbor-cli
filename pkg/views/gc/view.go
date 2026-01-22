package gc

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
	{Title: "ID", Width: 10},
	{Title: "Status", Width: 15},
	{Title: "Dry Run", Width: 10},
	{Title: "Creation Time", Width: 25},
	{Title: "Update Time", Width: 25},
}

func ListGC(history []*models.GCHistory) {
	var rows []table.Row
	for _, job := range history {
		creationTime, _ := utils.FormatCreatedTime(job.CreationTime.String())
		updateTime, _ := utils.FormatCreatedTime(job.UpdateTime.String())

		// Note: JobParameters is usually a JSON string. For simplicity we display it as is or handle parsing if needed.
		// Usually contains {"dry_run": true/false}

		rows = append(rows, table.Row{
			strconv.FormatInt(job.ID, 10),
			job.JobStatus,
			job.JobParameters,
			creationTime,
			updateTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
