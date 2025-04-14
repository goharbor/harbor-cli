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
package health

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/health"
	"github.com/goharbor/harbor-cli/pkg/views"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Component", Width: tablelist.WidthL},
	{Title: "Status", Width: tablelist.WidthXXL},
}

func PrintHealthStatus(status *health.GetHealthOK) {
	var rows []table.Row
	fmt.Printf("Harbor Health Status:: %s\n", styleStatus(status.Payload.Status))
	for _, component := range status.Payload.Components {
		rows = append(rows, table.Row{
			component.Name,
			styleStatus(component.Status),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func styleStatus(status string) string {
	if status == "healthy" {
		return views.GreenStyle.Render(status)
	}
	return views.RedStyle.Render(status)
}
