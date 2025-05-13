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
package metadata

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/table"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

func DisplayScannerMetadata(md *models.ScannerAdapterMetadata) {
	infoTable := buildInfoTable(md)
	capabilityTables := buildCapabilityTables(md)
	propTable := buildPropertyTable(md)

	fmt.Println("[Scanner Info]")
	fmt.Println(infoTable.View())

	fmt.Println("[Capabilities]")
	for _, capTable := range capabilityTables {
		fmt.Println(capTable.View())
	}

	fmt.Println("[Properties]")
	fmt.Println(propTable.View())
}

func buildInfoTable(md *models.ScannerAdapterMetadata) tablelist.Model {
	cols := []table.Column{
		{Title: "Key", Width: tablelist.WidthL},
		{Title: "Value", Width: tablelist.WidthXL},
	}
	rows := []table.Row{
		{"Name", md.Scanner.Name},
		{"Vendor", md.Scanner.Vendor},
		{"Version", md.Scanner.Version},
	}
	return tablelist.NewModel(cols, rows, len(rows))
}

func buildCapabilityTables(md *models.ScannerAdapterMetadata) []tablelist.Model {
	var tables []tablelist.Model

	for i, cap := range md.Capabilities {
		cols := []table.Column{
			{Title: fmt.Sprintf("Capability #%d - Consumes", i+1), Width: tablelist.WidthXXL * 2},
			{Title: "Produces", Width: tablelist.WidthXXL * 2},
		}

		maxLen := slices.Max([]int{len(cap.ConsumesMimeTypes), len(cap.ProducesMimeTypes)})
		var rows []table.Row
		for j := 0; j < maxLen; j++ {
			consume, produce := "", ""
			if j < len(cap.ConsumesMimeTypes) {
				consume = cap.ConsumesMimeTypes[j]
			}
			if j < len(cap.ProducesMimeTypes) {
				produce = cap.ProducesMimeTypes[j]
			}
			rows = append(rows, table.Row{consume, produce})
		}

		tables = append(tables, tablelist.NewModel(cols, rows, len(rows)))
	}
	return tables
}

func buildPropertyTable(md *models.ScannerAdapterMetadata) tablelist.Model {
	cols := []table.Column{
		{Title: "Property", Width: tablelist.WidthXXL * 2},
		{Title: "Value", Width: tablelist.WidthXXL * 2},
	}
	var rows []table.Row
	for k, v := range md.Properties {
		rows = append(rows, table.Row{k, v})
	}
	return tablelist.NewModel(cols, rows, len(rows))
}
