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
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	config "github.com/goharbor/harbor-cli/pkg/config/robot"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablegrid"
)

func NewRobotPermissionsGrid(kind string) *tablegrid.TableGrid {
	const (
		WidthResource = 20
		WidthAction   = 16
	)

	var resourceNames []string
	var columnLabels []string

	// Define resources, columns, and actions based on kind
	if kind == "project" {
		resourceNames = []string{
			"Accessory", "Artifact", "Artifact Addition", "Artifact Label",
			"Export CVE", "Immutable Tag", "Label", "Log", "Member",
			"Metadata", "Notification Policy", "Preheat Policy",
			"Project", "Quota", "Repository", "Robot", "SBOM",
			"Scan", "Scanner", "Tag", "Tag Retention",
		}
		columnLabels = []string{
			"Resource", "Create", "Delete",
			"List", "Pull", "Push", "Read",
			"Stop", "Update",
		}
	} else if kind == "system" {
		resourceNames = []string{
			"Audit Log", "Catalog", "Garbage Collection", "JobService Monitor",
			"Label", "LDAP User", "Preheat Instance", "Project", "Purge Audit",
			"Quota", "Registry", "Replication", "Replication Adapter", "Replication Policy",
			"Robot", "Scan All", "Scanner", "Security Hub", "System Volumes",
			"User", "User Group",
		}
		columnLabels = []string{
			"Resource", "Create", "Delete",
			"List", "Read", "Stop", "Update",
		}
	} else {
		fmt.Printf("Invalid kind specified: %s, expected 'system' or 'project'\n", kind)
		os.Exit(1)
	}

	actions := make([]string, len(columnLabels)-1)
	for i := 1; i < len(columnLabels); i++ {
		actions[i-1] = strings.ToLower(columnLabels[i])
	}

	// System Robot Permissions Grid
	// The disabled map corresponds to the permissions for each resource and action in the harbor UI.
	// This shows which actions are available (✓) or unavailable (✗) for each resource:
	// Resource              | Create | Delete | List | Read | Stop | Update |
	// ----------------------|--------|--------|------|------|------|--------|
	// Audit Log             |   ✗    |   ✗    |  ✓   |  ✗   |  ✗   |   ✗    |
	// Catalog               |   ✗    |   ✗    |  ✗   |  ✓   |  ✗   |   ✗    |
	// Garbage Collection    |   ✓    |   ✗    |  ✓   |  ✓   |  ✓   |   ✓    |
	// JobService Monitor    |   ✗    |   ✗    |  ✓   |  ✗   |  ✓   |   ✗    |
	// Label                 |   ✓    |   ✓    |  ✗   |  ✓   |  ✗   |   ✓    |
	// LDAP User             |   ✓    |   ✗    |  ✓   |  ✗   |  ✗   |   ✗    |
	// Preheat Instance      |   ✓    |   ✓    |  ✓   |  ✓   |  ✗   |   ✓    |
	// Project               |   ✓    |   ✗    |  ✓   |  ✗   |  ✗   |   ✗    |
	// Purge Audit           |   ✓    |   ✗    |  ✓   |  ✓   |  ✓   |   ✓    |
	// Quota                 |   ✗    |   ✗    |  ✓   |  ✓   |  ✗   |   ✓    |
	// Registry              |   ✓    |   ✓    |  ✓   |  ✓   |  ✗   |   ✓    |
	// Replication           |   ✓    |   ✗    |  ✓   |  ✓   |  ✗   |   ✗    |
	// Replication Adapter   |   ✗    |   ✗    |  ✓   |  ✗   |  ✗   |   ✗    |
	// Replication Policy    |   ✓    |   ✓    |  ✓   |  ✓   |  ✗   |   ✓    |
	// Robot                 |   ✓    |   ✓    |  ✓   |  ✓   |  ✗   |   ✗    |
	// Scan All              |   ✓    |   ✗    |  ✗   |  ✓   |  ✓   |   ✓    |
	// Scanner               |   ✓    |   ✓    |  ✓   |  ✓   |  ✗   |   ✓    |
	// Security Hub          |   ✗    |   ✗    |  ✓   |  ✓   |  ✗   |   ✗    |
	// System Volumes        |   ✗    |   ✗    |  ✗   |  ✓   |  ✗   |   ✗    |
	// User                  |   ✓    |   ✓    |  ✓   |  ✓   |  ✗   |   ✓    |
	// User Group            |   ✓    |   ✓    |  ✓   |  ✓   |  ✗   |   ✓    |
	//
	// -----------------------------------------------------------------------
	//
	// Project Robot Permissions Grid
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
	// This corresponds to the following. the map is autogenerated based on the available permissions
	// and the resource names. The indices in the map correspond to the column indices in the table.

	// Set up column widths
	columnWidths := []int{WidthResource}
	for range len(columnLabels) - 1 {
		columnWidths = append(columnWidths, WidthAction)
	}

	// Get available permissions from API
	perms, err := config.GetAllAvailablePermissions()
	if err != nil {
		fmt.Printf("Error fetching available permissions: %v\n", err)
		os.Exit(1)
	}

	var availablePerms map[string][]string
	if kind == "project" {
		availablePerms = perms.Project
	} else if kind == "system" {
		availablePerms = perms.System
	} else {
		fmt.Printf("invalid kind specified: %s, expected 'system' or 'project'\n", kind)
		os.Exit(1)
	}

	// Map display names to API resource keys
	resourceKeyMap := make(map[string]string)
	for _, displayName := range resourceNames {
		kebabName := utils.ToKebabCase(displayName)
		resourceKeyMap[displayName] = kebabName
	}

	// Create ordered list of resources that exist in API
	orderedResources := []string{}
	processedKeys := make(map[string]bool)

	// First add resources in our predefined order that exist in API
	for _, displayName := range resourceNames {
		kebabName := resourceKeyMap[displayName]
		if _, exists := availablePerms[kebabName]; exists {
			orderedResources = append(orderedResources, displayName)
			processedKeys[kebabName] = true
		}
	}

	// Create disabled map for UI grid
	disabled := make(map[int]map[int]bool)
	for rowIdx := range orderedResources {
		disabled[rowIdx] = make(map[int]bool)

		// For each action column
		for colIdx, action := range actions {
			resourceKey := resourceKeyMap[orderedResources[rowIdx]]
			validActions := availablePerms[resourceKey]

			// Disable action if it's not available for this resource
			disabled[rowIdx][colIdx+1] = !slices.Contains(validActions, action)
		}
	}

	// Create the table grid
	icons := &tablegrid.Icons{
		Selected:   "✅",
		Unselected: "❌",
		Empty:      " ",
	}

	return tablegrid.New(tablegrid.Config{
		RowLabels:    orderedResources,
		ColLabels:    columnLabels,
		Disabled:     disabled,
		ColumnWidths: columnWidths,
		Icons:        icons,
		Footer:       "\n ↑/↓ move row • ←/→ move col • space/enter to toggle • ⌃A toggle row • q to cancel\n",
	})
}

func ListPermissions(perms *models.Permissions, kind string, ch chan<- []models.Permission) {
	grid := NewRobotPermissionsGrid(kind)
	_, err := tea.NewProgram(grid, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	data := grid.GetData()
	selectedPerms := []models.Permission{}

	for rowIdx, displayName := range grid.RowLabels {
		kebabResource := utils.ToKebabCase(displayName)

		for colIdx := 0; colIdx < len(grid.ColLabels)-1; colIdx++ {
			if data[rowIdx][colIdx] {
				action := strings.ToLower(grid.ColLabels[colIdx+1])
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
