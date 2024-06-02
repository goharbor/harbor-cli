package list

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
	{Title: "Name", Width: 24},
	{Title: "Artifacts", Width: 12},
	{Title: "Pulls", Width: 12},
	{Title: "Last Modified Time", Width: 30},
}

func ListRepositories(repos []*models.Repository) {
	var rows []table.Row
	for _, repo := range repos {

		createdTime, _ := utils.FormatCreatedTime(repo.UpdateTime.String())
		rows = append(rows, table.Row{
			repo.Name,
			fmt.Sprintf("%d", repo.ArtifactCount),
			strconv.FormatInt(repo.PullCount, 10),
			createdTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
