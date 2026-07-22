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
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	config "github.com/goharbor/harbor-cli/pkg/config/robot"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/multiselectv2"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablegridv2"
)

type PermissionSelectResult struct {
	Permissions []models.Permission
	Err         error
}

func NewRobotPermissionsGridConfig(kind string) (tablegridv2.Config, error) {
	const (
		WidthResource = 20
		WidthAction   = 16
	)

	// Get available permissions from API
	perms, err := config.GetAllAvailablePermissions()
	if err != nil {
		return tablegridv2.Config{}, fmt.Errorf("error fetching available permissions: %w", err)
	}

	var availablePerms map[string][]string
	switch kind {
	case "project":
		availablePerms = perms.Project
	case "system":
		availablePerms = perms.System
	default:
		return tablegridv2.Config{}, fmt.Errorf("invalid kind specified: %s, expected 'system' or 'project'", kind)
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
	icons := &tablegridv2.Icons{
		Selected:   "✅",
		Unselected: "❌",
		Empty:      " ",
	}

	return tablegridv2.Config{
		RowLabels:    orderedResources,
		ColLabels:    columnLabels,
		Disabled:     disabled,
		ColumnWidths: columnWidths,
		Icons:        icons,
		Footer:       "\n ↑/↓ move row • ←/→ move col • space/enter to toggle • ⌃A toggle row • q to cancel\n",
	}, nil
}

func ListPermissions(kind string) PermissionSelectResult {
	switch kind {
	case "system":
		return listPermissionsMultiSelect(kind)
	case "project":
		return listPermissionsGrid(kind)
	default:
		return PermissionSelectResult{
			Permissions: nil,
			Err:         fmt.Errorf("invalid permission selection kind: %s", kind),
		}
	}
}

func listPermissionsGrid(kind string) PermissionSelectResult {
	model := tablegridv2.NewModel(
		fmt.Sprintf("Loading %s permissions...", kind),
		func() (tablegridv2.Config, error) {
			return NewRobotPermissionsGridConfig(kind)
		},
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return PermissionSelectResult{
			Permissions: nil,
			Err:         fmt.Errorf("error creating permissions grid: %w", err),
		}
	}

	loadedModel, ok := finalModel.(tablegridv2.Model)
	if !ok {
		return PermissionSelectResult{
			Permissions: nil,
			Err:         errors.New("unexpected permission grid model result"),
		}
	}
	if loadedModel.Err != nil {
		return PermissionSelectResult{
			Permissions: nil,
			Err:         loadedModel.Err,
		}
	}
	if len(loadedModel.RowLabels) == 0 || len(loadedModel.ColLabels) == 0 {
		return PermissionSelectResult{
			Permissions: nil,
			Err:         errors.New("permission grid was not initialized"),
		}
	}

	data := loadedModel.Data
	selectedPerms := make([]models.Permission, 0)

	for rowIdx, displayName := range loadedModel.RowLabels {
		kebabResource := utils.ToKebabCase(displayName)

		for colIdx := 0; colIdx < len(loadedModel.ColLabels)-1; colIdx++ {
			if data[rowIdx][colIdx] {
				action := strings.ToLower(loadedModel.ColLabels[colIdx+1])
				selectedPerms = append(selectedPerms, models.Permission{
					Resource: kebabResource,
					Action:   action,
				})
			}
		}
	}

	return PermissionSelectResult{
		Permissions: selectedPerms,
		Err:         nil,
	}
}

func listPermissionsMultiSelect(kind string) PermissionSelectResult {
	selected := []models.Permission{}
	model := multiselectv2.NewModel(
		&selected,
		fmt.Sprintf("Loading %s permissions...", kind),
		func() ([]models.Permission, error) {
			response, err := api.GetPermissions()
			if err != nil {
				return nil, err
			}

			var permissionSet []*models.Permission
			switch kind {
			case "system":
				permissionSet = response.Payload.System
			case "project":
				permissionSet = response.Payload.Project
			default:
				return nil, fmt.Errorf("invalid permission selection kind: %s", kind)
			}

			choices := make([]models.Permission, 0, len(permissionSet))
			for _, permission := range permissionSet {
				choices = append(choices, *permission)
			}

			if len(choices) == 0 {
				return nil, fmt.Errorf("no %s permissions found", kind)
			}

			return choices, nil
		},
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return PermissionSelectResult{
			Permissions: nil,
			Err:         fmt.Errorf("error creating permissions selector: %w", err),
		}
	}

	loadedModel, ok := finalModel.(multiselectv2.Model)
	if !ok {
		return PermissionSelectResult{
			Permissions: nil,
			Err:         errors.New("unexpected permission multiselect model result"),
		}
	}
	if loadedModel.Err != nil {
		return PermissionSelectResult{
			Permissions: nil,
			Err:         loadedModel.Err,
		}
	}

	return PermissionSelectResult{
		Permissions: selected,
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
