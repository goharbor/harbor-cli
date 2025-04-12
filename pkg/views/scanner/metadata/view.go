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
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

func DisplayScannerMetadata(md *models.ScannerAdapterMetadata) {
	infoCols := []table.Column{
		{Title: "Key", Width: 15},
		{Title: "Value", Width: 50},
	}
	infoRows := []table.Row{
		{"Name", md.Scanner.Name},
		{"Vendor", md.Scanner.Vendor},
		{"Version", md.Scanner.Version},
	}

	capCols := []table.Column{
		{Title: "Consumes", Width: 45},
		{Title: "Produces", Width: 45},
	}
	var capRows []table.Row
	for _, cap := range md.Capabilities {
		consumes := joinOrNone(cap.ConsumesMimeTypes)
		produces := joinOrNone(cap.ProducesMimeTypes)
		capRows = append(capRows, table.Row{consumes, produces})
	}

	propCols := []table.Column{
		{Title: "Property", Width: 50},
		{Title: "Value", Width: 50},
	}
	var propRows []table.Row
	for k, v := range md.Properties {
		propRows = append(propRows, table.Row{k, v})
	}

	m := metadataModel{
		infoTable: tablelist.NewModel(infoCols, infoRows, len(infoRows)),
		capTable:  tablelist.NewModel(capCols, capRows, len(capRows)),
		propTable: tablelist.NewModel(propCols, propRows, len(propRows)),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type metadataModel struct {
	infoTable tablelist.Model
	capTable  tablelist.Model
	propTable tablelist.Model
}

func (m metadataModel) Init() tea.Cmd {
	return nil
}

func (m metadataModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m metadataModel) View() string {
	return fmt.Sprintf(
		"\n[Scanner Info]\n%s\n[Capabilities]\n%s\n[Properties]\n%s\n",
		m.infoTable.View(),
		m.capTable.View(),
		m.propTable.View(),
	)
}

func joinOrNone(list []string) string {
	if len(list) == 0 {
		return "None"
	}
	return fmt.Sprintf("%v", list)
}
