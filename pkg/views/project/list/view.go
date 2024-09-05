// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
