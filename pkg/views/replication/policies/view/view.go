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
	"strconv"
	"strings"

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

var filterColumns = []table.Column{
	{Title: "Type", Width: tablelist.WidthM},
	{Title: "Decoration", Width: tablelist.WidthM},
	{Title: "Value", Width: tablelist.WidthXL},
}

var order = []string{
	"ID",
	"Name",
	"Enabled",
	"Source",
	"Destination",
	"Creation Time",
	"Last Modified",
	"Description",
	"Trigger Type",
	"Override",
	"Replicate Deletion",
	"Copy By Chunk",
	"Speed",
}

func ViewPolicy(rpolicy *models.ReplicationPolicy) {
	createdTime, _ := utils.FormatCreatedTime(rpolicy.CreationTime.String())
	modifledTime, _ := utils.FormatCreatedTime(rpolicy.UpdateTime.String())
	policyMap := map[string]string{
		"ID":                 strconv.FormatInt(rpolicy.ID, 10),
		"Name":               rpolicy.Name,
		"Source":             getRegistryName(rpolicy.SrcRegistry),
		"Destination":        getRegistryName(rpolicy.DestRegistry),
		"Trigger Type":       rpolicy.Trigger.Type,
		"Override":           strconv.FormatBool(rpolicy.Override),
		"Enabled":            strconv.FormatBool(rpolicy.Enabled),
		"Creation Time":      createdTime,
		"Last Modified":      modifledTime,
		"Description":        rpolicy.Description,
		"Replicate Deletion": strconv.FormatBool(rpolicy.ReplicateDeletion),
		"Copy By Chunk":      strconv.FormatBool(*rpolicy.CopyByChunk),
		"Speed":              strconv.FormatInt(int64(*rpolicy.Speed), 10) + " B/s",
	}

	var rows []table.Row
	for _, key := range order {
		rows = append(rows, table.Row{
			key,
			policyMap[key],
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if len(rpolicy.Filters) > 0 {
		fmt.Println("\nFilters:")
		showFiltersTable(rpolicy.Filters)
	} else {
		fmt.Println("\nNo filters configured")
	}
}

func showFiltersTable(filters []*models.ReplicationFilter) {
	var filterRows []table.Row

	for _, filter := range filters {
		decoration := filter.Decoration
		if decoration == "" {
			decoration = "N/A"
		}

		value := formatFilterValue(filter.Value)

		filterRows = append(filterRows, table.Row{
			filter.Type,
			decoration,
			value,
		})
	}

	m := tablelist.NewModel(filterColumns, filterRows, len(filterRows))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func formatFilterValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ", ")
	case []interface{}:
		var strValues []string
		for _, item := range v {
			strValues = append(strValues, fmt.Sprintf("%v", item))
		}
		return strings.Join(strValues, ", ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getRegistryName(registry *models.Registry) string {
	if registry == nil {
		return "Local"
	}
	return registry.Name
}

func filtersToString(filters []*models.ReplicationFilter) string {
	if len(filters) == 0 {
		return "No filters"
	}

	var filterStrings []string
	for _, filter := range filters {
		filterStrings = append(filterStrings, fmt.Sprintf("%s:%s", filter.Type, filter.Value))
	}
	return fmt.Sprintf("[%s]", joinWithComma(filterStrings...))

}

func joinWithComma(elements ...string) string {
	if len(elements) == 0 {
		return ""
	}
	result := elements[0]
	for _, elem := range elements[1:] {
		result += ", " + elem
	}
	return result
}
