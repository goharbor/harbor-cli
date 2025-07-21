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
	{Title: "User", Width: 18},
	{Title: "Resource", Width: 18},
	{Title: "Resource Type", Width: 14},
	{Title: "Operation", Width: 10},
	{Title: "Time", Width: 16},
}

func ListLogs(logs []*models.AuditLogExt) {
	var rows []table.Row
	for _, log := range logs {
		operationTime, _ := utils.FormatCreatedTime(log.OpTime.String())
		rows = append(rows, table.Row{
			log.Username,
			log.Resource,
			log.ResourceType,
			log.Operation,
			operationTime,
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
