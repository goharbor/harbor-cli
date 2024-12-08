package search

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Repository Name", Width: 30},
	{Title: "Project Id", Width: 11},
	{Title: "Project Name", Width: 16},
	{Title: "Access Level", Width: 12},
	{Title: "Artifact Count", Width: 14},
	{Title: "Pull Count", Width: 12},
}

func SearchRepositories(repos []*models.SearchRepository) {
	var rows []table.Row
	for _, repo := range repos {
		accessLevel := "public"
		if !repo.ProjectPublic {
			accessLevel = "private"
		}
		rows = append(rows, table.Row{
			repo.RepositoryName,
			fmt.Sprintf("%d", repo.ProjectID),
			repo.ProjectName,
			accessLevel,
			fmt.Sprintf("%d", repo.ArtifactCount),
			strconv.FormatInt(repo.PullCount, 10),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
