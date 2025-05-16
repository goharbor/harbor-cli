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
	{Title: "Name", Width: tablelist.WidthM},
	{Title: "Description", Width: tablelist.WidthXXL},
	{Title: "Disabled", Width: tablelist.WidthM},
	{Title: "URL", Width: tablelist.WidthXXL},
	{Title: "Default", Width: tablelist.WidthS},
	{Title: "Internal", Width: tablelist.WidthM},
	{Title: "Created", Width: tablelist.WidthM},
}

func ListScanners(scanners []*models.ScannerRegistration) {
	var rows []table.Row
	for _, s := range scanners {
		createdAt, err := utils.FormatCreatedTime(s.CreateTime.String())
		if err != nil {
			fmt.Println("Error formatting created time:", err)
			os.Exit(1)
		}
		rows = append(rows, table.Row{
			s.Name,
			s.Description,
			fmt.Sprintf("%v", *s.Disabled),
			fmt.Sprintf("%v", s.URL),
			fmt.Sprintf("%v", *s.IsDefault),
			fmt.Sprintf("%v", *s.UseInternalAddr),
			fmt.Sprintf("%v", createdAt),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
