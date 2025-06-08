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

package robot

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablegrid"
)

const (
	WidthResource = 20
	WidthAction   = 16
)

var resourceNames = []string{
	"Accessory",
	"Artifact",
	"Artifact Addition",
	"Artifact Label",
	"Export CVE",
	"Immutable Tag",
	"Label",
	"Log",
	"Member",
	"Metadata",
	"Notification Policy",
	"Preheat Policy",
	"Project",
	"Quota",
	"Repository",
	"Robot Account",
	"SBOM",
	"Scan",
	"Scanner",
	"Tag",
	"Tag Retention",
}

var columnLabels = []string{
	"Resource", "Create", "Delete", "List", "Pull", "Push", "Read", "Stop", "Update",
}

func NewRobotPermissionsGrid() *tablegrid.TableGrid {
	columnWidths := []int{WidthResource}
	for range len(columnLabels) - 1 {
		columnWidths = append(columnWidths, WidthAction)
	}

	// The disabled map corresponds to the permissions for each resource and action in the harbor UI.
	// This shows which actions are available (✓) or unavailable (✗) for each resource:
	//
	// Resource           | Create | Delete | List | Pull | Push | Read | Stop | Update |
	// -------------------|--------|--------|------|------|------|------|------|--------|
	// Accessory          |   ✗    |   ✗    |  ✓   |  ✗   |  ✗   |  ✗   |  ✗   |   ✗    |
	// Artifact           |   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✓   |  ✗   |   ✗    |
	// Artifact Addition  |   ✗    |   ✗    |  ✗   |  ✗   |  ✗   |  ✓   |  ✗   |   ✗    |
	// Artifact Label     |   ✓    |   ✓    |  ✗   |  ✗   |  ✗   |  ✗   |  ✗   |   ✗    |
	// Export CVE         |   ✓    |   ✗    |  ✗   |  ✗   |  ✗   |  ✓   |  ✗   |   ✗    |
	// Immutable Tag      |   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✗   |  ✗   |   ✓    |
	// Label              |   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✓   |  ✗   |   ✓    |
	// Log                |   ✗    |   ✗    |  ✓   |  ✗   |  ✗   |  ✗   |  ✗   |   ✗    |
	// Member             |   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✓   |  ✗   |   ✓    |
	// Metadata           |   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✓   |  ✗   |   ✓    |
	// Notification Policy|   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✓   |  ✗   |   ✓    |
	// Preheat Policy     |   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✓   |  ✗   |   ✓    |
	// Project            |   ✗    |   ✗    |  ✗   |  ✗   |  ✗   |  ✓   |  ✓   |   ✗    |
	// Quota              |   ✗    |   ✗    |  ✗   |  ✗   |  ✗   |  ✓   |  ✗   |   ✗    |
	// Repository         |   ✗    |   ✓    |  ✓   |  ✓   |  ✓   |  ✓   |  ✗   |   ✓    |
	// Robot Account      |   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✓   |  ✗   |   ✗    |
	// SBOM               |   ✓    |   ✗    |  ✗   |  ✗   |  ✗   |  ✓   |  ✓   |   ✗    |
	// Scan               |   ✓    |   ✗    |  ✗   |  ✗   |  ✗   |  ✓   |  ✓   |   ✗    |
	// Scanner            |   ✓    |   ✗    |  ✗   |  ✗   |  ✗   |  ✓   |  ✗   |   ✗    |
	// Tag                |   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✗   |  ✗   |   ✗    |
	// Tag Retention      |   ✓    |   ✓    |  ✓   |  ✗   |  ✗   |  ✓   |  ✗   |   ✓    |
	//
	// Where:
	// - ✓ means the action is available for selection (disabled[row][col] = false)
	// - ✗ means the action is unavailable (disabled[row][col] = true)
	// - Indices in the map: 1=Create, 2=Delete, 3=List, 4=Pull, 5=Push, 6=Read, 7=Stop, 8=Update
	disabled := map[int]map[int]bool{
		0:  {1: true, 2: true, 3: false, 4: true, 5: true, 6: true, 7: true, 8: true},
		1:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: true},
		2:  {1: true, 2: true, 3: true, 4: true, 5: true, 6: false, 7: true, 8: true},
		3:  {1: false, 2: false, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true},
		4:  {1: false, 2: true, 3: true, 4: true, 5: true, 6: false, 7: true, 8: true},
		5:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: true, 7: true, 8: false},
		6:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		7:  {1: true, 2: true, 3: false, 4: true, 5: true, 6: true, 7: true, 8: true},
		8:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		9:  {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		10: {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		11: {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
		12: {1: true, 2: true, 3: true, 4: true, 5: true, 6: false, 7: false, 8: true},
		13: {1: true, 2: true, 3: true, 4: true, 5: true, 6: false, 7: true, 8: true},
		14: {1: true, 2: false, 3: false, 4: false, 5: false, 6: false, 7: true, 8: false},
		15: {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: true},
		16: {1: false, 2: true, 3: true, 4: true, 5: true, 6: false, 7: false, 8: true},
		17: {1: false, 2: true, 3: true, 4: true, 5: true, 6: false, 7: false, 8: true},
		18: {1: false, 2: true, 3: true, 4: true, 5: true, 6: false, 7: true, 8: true},
		19: {1: false, 2: false, 3: false, 4: true, 5: true, 6: true, 7: true, 8: true},
		20: {1: false, 2: false, 3: false, 4: true, 5: true, 6: false, 7: true, 8: false},
	}

	icons := &tablegrid.Icons{
		Selected:   "✅",
		Unselected: "❌",
		Empty:      " ",
	}

	return tablegrid.New(tablegrid.Config{
		RowLabels:    resourceNames,
		ColLabels:    columnLabels,
		Disabled:     disabled,
		ColumnWidths: columnWidths,
		Icons:        icons,
		Footer:       "\n ↑/↓ move row • ←/→ move col • space/enter to toggle • ⌃A toggle row • q to cancel\n",
	})
}

func ListPermissions(perms *models.Permissions, ch chan<- []models.Permission) {
	grid := NewRobotPermissionsGrid()
	_, err := tea.NewProgram(grid, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	data := grid.GetData()
	selectedPerms := []models.Permission{}

	for rowIdx, res := range resourceNames {
		kebabResource := strings.ToLower(res)
		kebabResource = strings.ReplaceAll(kebabResource, " ", "-")

		for colIdx := 0; colIdx < len(columnLabels)-1; colIdx++ {
			if data[rowIdx][colIdx] {
				action := strings.ToLower(columnLabels[colIdx+1])
				selectedPerms = append(selectedPerms, models.Permission{
					Resource: kebabResource,
					Action:   action,
				})
			}
		}
	}

	ch <- selectedPerms
}

func ListRobot(robots []*models.Robot, choice chan<- int64) {
	itemsList := make([]list.Item, len(robots))
	items := map[string]int64{}
	for i, r := range robots {
		items[r.Name] = r.ID
		itemsList[i] = selection.Item(r.Name)
	}
	m := selection.NewModel(itemsList, "Robot")
	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if pm, ok := p.(selection.Model); ok {
		choice <- items[pm.Choice]
	}
}
