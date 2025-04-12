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
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var columns = []table.Column{
	{Title: "Name", Width: 10},
	{Title: "Description", Width: 30},
	{Title: "Default", Width: 8},
	{Title: "URL", Width: 30},
	{Title: "Internal Addr", Width: 14},
	{Title: "Created At", Width: 20},
	{Title: "Disabled", Width: 20},
}

func formatTime(raw string) string {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}
	return t.Format("02 Jan 2006, 15:04")
}

func ListScanners(scanners []*models.ScannerRegistration) {
	var rows []table.Row
	for _, s := range scanners {
		rows = append(rows, table.Row{
			s.Name,
			s.Description,
			fmt.Sprintf("%v", *s.IsDefault),
			fmt.Sprintf("%v", s.URL),
			fmt.Sprintf("%v", *s.UseInternalAddr),
			formatTime(s.CreateTime.String()),
			fmt.Sprintf("%v", *s.Disabled),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
