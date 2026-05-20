// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package summary

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Property", Width: tablelist.WidthL},
	{Title: "Value", Width: tablelist.Width3XL},
}


var order = []string{
	"Project Name",
	"Project ID",
	"Visibility",
	"Repositories",
	"Project Admin Count",
	"Maintainer Count",
	"Developer Count",
	"Guest Count",
	"Limited Guest Count",
}

func ViewProjectSummary(project *models.Project, summary *models.ProjectSummary) error {
	summaryMap := map[string]string{
		"Project Name":        project.Name,
		"Project ID":          strconv.FormatInt(int64(project.ProjectID), 10),
		"Visibility":          project.Metadata.Public,
		"Repositories":        strconv.FormatInt(summary.RepoCount, 10),
		"Project Admin Count": strconv.FormatInt(summary.ProjectAdminCount, 10),
		"Maintainer Count":    strconv.FormatInt(summary.MaintainerCount, 10),
		"Developer Count":     strconv.FormatInt(summary.DeveloperCount, 10),
		"Guest Count":         strconv.FormatInt(summary.GuestCount, 10),
		"Limited Guest Count": strconv.FormatInt(summary.LimitedGuestCount, 10),
	}

	var rows []table.Row
	for _, key := range order {
		rows = append(rows, table.Row{key, summaryMap[key]})
	}

	if summary.Quota != nil {
		for resource, hardValue := range summary.Quota.Hard {
			usedValue := summary.Quota.Used[resource]
			label := fmt.Sprintf("Quota: %s", resource)
			val := ""
			if resource == "storage" {
				val = fmt.Sprintf("%v / %v", utils.FormatSize(usedValue), utils.FormatSize(hardValue))
			} else {
				val = fmt.Sprintf("%v / %v", utils.FormatCount(usedValue), utils.FormatCount(hardValue))
			}
			rows = append(rows, table.Row{label, val})
		}
	}

	if summary.Registry != nil {
		rows = append(rows, table.Row{"Registry ID", strconv.FormatInt(summary.Registry.ID, 10)})
		rows = append(rows, table.Row{"Registry Name", summary.Registry.Name})
		rows = append(rows, table.Row{"Registry URL", summary.Registry.URL})
		rows = append(rows, table.Row{"Registry Type", summary.Registry.Type})
		if summary.Registry.Credential != nil {
			rows = append(rows, table.Row{"Registry Credential Type", summary.Registry.Credential.Type})
		}
		rows = append(rows, table.Row{"Registry Description", summary.Registry.Description})
		rows = append(rows, table.Row{"Registry Insecure", strconv.FormatBool(summary.Registry.Insecure)})
		rows = append(rows, table.Row{"Registry Status", summary.Registry.Status})
		rows = append(rows, table.Row{"Registry Updated", summary.Registry.UpdateTime.String()})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		return fmt.Errorf("error running program: %v", err)
	}
	return nil
}
