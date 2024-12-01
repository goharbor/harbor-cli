package view

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
	{Title: "Name", Width: 20},
	{Title: "ID", Width: 10},
	{Title: "Project ID", Width: 10},
	{Title: "Artifacts", Width: 10},
	{Title: "Pulls", Width: 5},
	{Title: "Creation Time", Width: 20},
	{Title: "Last Modified Time", Width: 20},
	{Title: "Description", Width: 20},
}

func ViewRepository(repo *models.Repository) {
	var rows []table.Row

	createdTime, _ := utils.FormatCreatedTime(repo.CreationTime.String())
	modifledTime, _ := utils.FormatCreatedTime(repo.UpdateTime.String())
	rows = append(rows, table.Row{
		repo.Name,
		fmt.Sprintf("%d", repo.ID),
		fmt.Sprintf("%d", repo.ProjectID),
		fmt.Sprintf("%d", repo.ArtifactCount),
		strconv.FormatInt(repo.PullCount, 10),
		createdTime,
		modifledTime,
		repo.Description,
	})

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
