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
	"strconv"
	"strings"

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

var eventTypeColumns = []table.Column{
	{Title: "INDEX", Width: tablelist.WidthS},
	{Title: "EVENT_TYPE", Width: 40},
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

func ListAuditLogEventTypes(eventTypes []*models.AuditLogEventType, page, pageSize int64, total int, showPaginationSummary bool) {
	if len(eventTypes) == 0 {
		if showPaginationSummary {
			fmt.Println("No audit log event types found for the requested page.")
			return
		}
		fmt.Println("No audit log event types found.")
		return
	}

	startIndex := int64(1)
	if showPaginationSummary {
		startIndex = (page-1)*pageSize + 1
	}

	var rows []table.Row
	for i, eventType := range eventTypes {
		rows = append(rows, table.Row{
			strconv.FormatInt(startIndex+int64(i), 10),
			auditLogEventTypeName(eventType),
		})
	}

	m := tablelist.NewModel(eventTypeColumns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if showPaginationSummary {
		endIndex := startIndex + int64(len(eventTypes)) - 1
		fmt.Printf("\nShowing %d-%d of %d\n", startIndex, endIndex, total)
	}
}

func auditLogEventTypeName(eventType *models.AuditLogEventType) string {
	if eventType == nil {
		return "-"
	}

	name := strings.TrimSpace(eventType.EventType)
	if name == "" {
		return "-"
	}

	return name
}
