package list

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var column = []table.Column{
	{Title: "Public Project", Width: 14},
	{Title: "Public Repo", Width: 12},
	{Title: "Private Project", Width: 14},
	{Title: "Private Repo", Width: 12},
}

var columnWide = []table.Column{
	{Title: "Public Project", Width: 14},
	{Title: "Public Repo", Width: 12},
	{Title: "Private Project", Width: 14},
	{Title: "Private Repo", Width: 12},
	{Title: "Total Project Count", Width: 16},
	{Title: "Total Repo Count", Width: 16},
	{Title: "Total Storage Consumption", Width: 16},
}

func ListStatistics(stat *models.Statistic, wide bool) {
	var rows []table.Row
	var columns []table.Column

	if wide {
		columns = columnWide
		rows = append(rows, table.Row{
			strconv.FormatInt(stat.PublicProjectCount, 10),
			strconv.FormatInt(stat.PublicRepoCount, 10),
			strconv.FormatInt(stat.PrivateProjectCount, 10),
			strconv.FormatInt(stat.PrivateRepoCount, 10),
			strconv.FormatInt(stat.TotalProjectCount, 10),
			strconv.FormatInt(stat.TotalRepoCount, 10),
			strconv.FormatInt(stat.TotalStorageConsumption, 10),
		})
	} else {
		columns = column
		rows = append(rows, table.Row{
			strconv.FormatInt(stat.PublicProjectCount, 10),
			strconv.FormatInt(stat.PublicRepoCount, 10),
			strconv.FormatInt(stat.PrivateProjectCount, 10),
			strconv.FormatInt(stat.PrivateRepoCount, 10),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
