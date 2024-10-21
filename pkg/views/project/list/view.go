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
	{Title: "ID", Width: 6},
	{Title: "Project Name", Width: 12},
	{Title: "Access Level", Width: 12},
	{Title: "Type", Width: 12},
	{Title: "Repo Count", Width: 12},
	{Title: "Creation Time", Width: 30},
}

func ListProjects(projects []*models.Project) {
	var rows []table.Row
	for _, project := range projects {
		accessLevel := "public"
		if project.Metadata.Public != "true" {
			accessLevel = "private"
		}

		projectType := "project"

		if project.RegistryID != 0 {
			projectType = "proxy cache"
		}
		createdTime, _ := utils.FormatCreatedTime(project.CreationTime.String())
		rows = append(rows, table.Row{
			strconv.FormatInt(int64(project.ProjectID), 10), // ProjectID
			project.Name, // Project Name
			accessLevel,  // Access Level
			projectType,  // Type
			strconv.FormatInt(project.RepoCount, 10),
			createdTime, // Creation Time
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
