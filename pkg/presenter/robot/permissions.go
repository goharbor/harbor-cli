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
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	config "github.com/goharbor/harbor-cli/pkg/config/robot"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/multiselectv2"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablegridv2"
)

func SelectPermissions(kind string) ([]models.Permission, error) {
	switch kind {
	case "system":
		return selectPermissionsMulti(kind)
	case "project":
		return selectPermissionsGrid(kind)
	default:
		return nil, fmt.Errorf("invalid permission selection kind: %s", kind)
	}
}

func NewPermissionsGridConfig(kind string) (tablegridv2.Config, error) {
	const (
		widthResource = 20
		widthAction   = 16
	)

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

	resourceKeys := make([]string, 0, len(availablePerms))
	for key := range availablePerms {
		resourceKeys = append(resourceKeys, key)
	}
	slices.Sort(resourceKeys)

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

	columnLabels := append([]string{"Resource"}, allActions...)
	orderedResources := make([]string, len(resourceKeys))
	for i, kebabKey := range resourceKeys {
		orderedResources[i] = utils.FromKebabCase(kebabKey)
	}

	columnWidths := []int{widthResource}
	for range columnLabels[1:] {
		columnWidths = append(columnWidths, widthAction)
	}

	disabled := make(map[int]map[int]bool)
	for rowIdx, resourceKey := range resourceKeys {
		disabled[rowIdx] = make(map[int]bool)
		validActions := availablePerms[resourceKey]
		for colIdx, action := range allActions {
			disabled[rowIdx][colIdx+1] = !slices.Contains(validActions, action)
		}
	}

	return tablegridv2.Config{
		RowLabels:    orderedResources,
		ColLabels:    columnLabels,
		Disabled:     disabled,
		ColumnWidths: columnWidths,
		Icons: &tablegridv2.Icons{
			Selected:   "✅",
			Unselected: "❌",
			Empty:      " ",
		},
		Footer: "\n ↑/↓ move row • ←/→ move col • space/enter to toggle • ⌃A toggle row • q to cancel\n",
	}, nil
}

func selectPermissionsGrid(kind string) ([]models.Permission, error) {
	model := tablegridv2.NewModel(
		fmt.Sprintf("Loading %s permissions...", kind),
		func() (tablegridv2.Config, error) {
			return NewPermissionsGridConfig(kind)
		},
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return nil, fmt.Errorf("error creating permissions grid: %w", err)
	}

	loadedModel, ok := finalModel.(tablegridv2.Model)
	if !ok {
		return nil, errors.New("unexpected permission grid model result")
	}
	if loadedModel.Err != nil {
		return nil, loadedModel.Err
	}
	if len(loadedModel.RowLabels) == 0 || len(loadedModel.ColLabels) == 0 {
		return nil, errors.New("permission grid was not initialized")
	}

	selectedPerms := make([]models.Permission, 0)
	for rowIdx, displayName := range loadedModel.RowLabels {
		kebabResource := utils.ToKebabCase(displayName)
		for colIdx := 0; colIdx < len(loadedModel.ColLabels)-1; colIdx++ {
			if loadedModel.Data[rowIdx][colIdx] {
				action := strings.ToLower(loadedModel.ColLabels[colIdx+1])
				selectedPerms = append(selectedPerms, models.Permission{
					Resource: kebabResource,
					Action:   action,
				})
			}
		}
	}

	return selectedPerms, nil
}

func selectPermissionsMulti(kind string) ([]models.Permission, error) {
	selected := []models.Permission{}
	model := multiselectv2.NewModel(
		&selected,
		fmt.Sprintf("Loading %s permissions...", kind),
		loadPermissionChoices(kind),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return nil, fmt.Errorf("error creating permissions selector: %w", err)
	}

	loadedModel, ok := finalModel.(multiselectv2.Model)
	if !ok {
		return nil, errors.New("unexpected permission multiselect model result")
	}
	if loadedModel.Err != nil {
		return nil, loadedModel.Err
	}

	return selected, nil
}

func loadPermissionChoices(kind string) multiselectv2.Loader {
	return func() ([]models.Permission, error) {
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
	}
}
