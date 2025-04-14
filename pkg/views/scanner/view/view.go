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
package view

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
	{Title: "Name", Width: 12},
	{Title: "UUID", Width: 38},
	{Title: "URL", Width: 30},
	{Title: "Default", Width: 8},
	{Title: "Disabled", Width: 9},
	{Title: "Skip Cert Verify", Width: 17},
	{Title: "Internal Addr", Width: 14},
	{Title: "Created At", Width: 18},
	{Title: "Updated At", Width: 18},
}

func ViewScanner(scanner *models.ScannerRegistration) {
	var rows []table.Row

	createdAt, _ := utils.FormatCreatedTime(scanner.CreateTime.String())
	updatedAt, _ := utils.FormatCreatedTime(scanner.UpdateTime.String())

	rows = append(rows, table.Row{
		scanner.Name,
		scanner.UUID,
		utils.FormatUrl(scanner.URL.String()),
		fmt.Sprintf("%v", *scanner.IsDefault),
		fmt.Sprintf("%v", *scanner.Disabled),
		fmt.Sprintf("%v", *scanner.SkipCertVerify),
		fmt.Sprintf("%v", *scanner.UseInternalAddr),
		createdAt,
		updatedAt,
	})

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
