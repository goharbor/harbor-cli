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

type PermissionSelectResult struct {
	Permissions []models.Permission
	Err         error
}

func NewRobotPermissionsGrid(kind string) (*tablegrid.TableGrid, error) {
	const (
		WidthResource = 20
		WidthAction   = 16
	)

	// Get available permissions from API
	perms, err := config.GetAllAvailablePermissions()
	if err != nil {
		return nil, fmt.Errorf("error fetching available permissions: %w", err)
	}

	var availablePerms map[string][]string
	switch kind {
	case "project":
		availablePerms = perms.Project
	case "system":
		availablePerms = perms.System
	default:
		return nil, fmt.Errorf("invalid kind specified: %s, expected 'system' or 'project'", kind)
	}

	// Extract and sort resource keys for deterministic ordering
	resourceKeys := make([]string, 0, len(availablePerms))
	for key := range availablePerms {
		resourceKeys = append(resourceKeys, key)
	}
	slices.Sort(resourceKeys)

	// Extract all unique actions across all resources, then sort
	actionSet := make(map[string]bool)
	for _, actions := range availablePerms {
		for _, action := range actions {
			actionSet[action] = true
		}
	}
	allActions := make([]string, 0, len(actionSet))
	for action := range actionSet {
		allActions = append(allActions, action)
	}
	slices.Sort(allActions)

	// Build column labels from sorted actions
	columnLabels := append([]string{"Resource"}, allActions...)

	// Convert kebab-case back to display names for row labels
	orderedResources := make([]string, len(resourceKeys))
	for i, kebabKey := range resourceKeys {
		orderedResources[i] = utils.FromKebabCase(kebabKey)
	}

	// Set up column widths
	columnWidths := []int{WidthResource}
	for range columnLabels[1:] {
		columnWidths = append(columnWidths, WidthAction)
	}

	// Create disabled map for UI grid
	disabled := make(map[int]map[int]bool)
	for rowIdx, resourceKey := range resourceKeys {
		disabled[rowIdx] = make(map[int]bool)
		validActions := availablePerms[resourceKey]
		// For each action column (skip first "Resource" column)
		for colIdx, action := range allActions {
			disabled[rowIdx][colIdx+1] = !slices.Contains(validActions, action)
		}
	}

	// Create the table grid
	icons := &tablegrid.Icons{
		Selected:   "✅",
		Unselected: "❌",
		Empty:      " ",
	}

	grid := tablegrid.New(tablegrid.Config{
		RowLabels:    orderedResources,
		ColLabels:    columnLabels,
		Disabled:     disabled,
		ColumnWidths: columnWidths,
		Icons:        icons,
		Footer:       "\n ↑/↓ move row • ←/→ move col • space/enter to toggle • ⌃A toggle row • q to cancel\n",
	})
	return grid, nil
}

func ListPermissions(perms *models.Permissions, kind string, ch chan<- PermissionSelectResult) {
	grid, err := NewRobotPermissionsGrid(kind)
	if err != nil {
		fmt.Println("error creating permissions grid:", err)
		ch <- PermissionSelectResult{
			Permissions: nil,
			Err:         err,
		}
		return
	}
	_, err = tea.NewProgram(grid, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("error creating permissions grid:", err)
		ch <- PermissionSelectResult{
			Permissions: nil,
			Err:         err,
		}
		return
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
	ch <- PermissionSelectResult{
		Permissions: selectedPerms,
		Err:         nil,
	}
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
