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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
)

// promptBasicInfo asks for name, description, and duration when they were
// not given as flags.
func promptBasicInfo(view *robotview.CreateView) error {
	if view.Name == "" || view.Duration == 0 {
		fmt.Println("Opening interactive form for robot creation")
		if err := robotview.FormCreate(view); err != nil {
			return fmt.Errorf("failed to read robot input: %v", err)
		}
	}

	if view.Duration == 0 {
		return fmt.Errorf("failed to create robot: duration cannot be 0")
	}
	return nil
}

// collectCreatePermissions gathers the permissions for a new robot, either
// all available ones or interactively.
func collectCreatePermissions(level Level, view *robotview.CreateView, all bool) (*PermissionSet, error) {
	permissions := NewPermissionSet()

	if level == LevelProject {
		projectPermissions, err := collectSingleProjectPermissions(all, "create")
		if err != nil {
			return nil, err
		}
		permissions.Project[view.ProjectName] = projectPermissions
		return permissions, nil
	}

	if err := collectSystemPermissions(all, permissions); err != nil {
		return nil, err
	}
	if err := collectProjectPermissions(view.ProjectName, all, permissions); err != nil {
		return nil, err
	}
	return permissions, nil
}

// collectSingleProjectPermissions covers project-level robots, which have
// exactly one permission scope.
func collectSingleProjectPermissions(all bool, action string) ([]models.Permission, error) {
	if all {
		return AllPermissions("project")
	}

	permissions, err := prompt.GetRobotPermissionsFromUser("project")
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %v", utils.ParseHarborErrorMsg(err))
	}
	if len(permissions) == 0 {
		return nil, errNoPermissions(action)
	}
	return permissions, nil
}

func collectSystemPermissions(all bool, permissions *PermissionSet) error {
	if all {
		systemPermissions, err := AllPermissions("system")
		if err != nil {
			return err
		}
		permissions.System = systemPermissions
		return nil
	}

	selected, err := prompt.GetRobotPermissionsFromUser("system")
	if err != nil {
		return fmt.Errorf("failed to get permissions: %v", utils.ParseHarborErrorMsg(err))
	}
	if len(selected) == 0 {
		return errNoPermissions("create")
	}
	permissions.System = selected
	return nil
}

func collectProjectPermissions(projectName string, all bool, permissions *PermissionSet) error {
	if projectName != "" {
		projectPermissions, err := collectSingleProjectPermissions(all, "create")
		if err != nil {
			return err
		}
		permissions.Project[projectName] = projectPermissions
		return nil
	}

	if all {
		// --all-permission without a project means system-level permissions only.
		return nil
	}

	mode, err := robotview.SelectCreatePermissionMode()
	if err != nil {
		return fmt.Errorf("error selecting permission mode: %v", err)
	}

	switch mode {
	case "list":
		return collectMultiProjectPermissions(permissions, false)
	case "per_project":
		return collectPerProjectPermissions(permissions, false)
	case "none":
		fmt.Println("Creating robot with system-level permissions only (no project-specific permissions)")
		return nil
	default:
		return fmt.Errorf("unknown permission mode: %s", mode)
	}
}

// collectMultiProjectPermissions applies one shared permission selection to
// several projects. With validate set, invalid project permissions are
// filtered out (the update flow).
func collectMultiProjectPermissions(permissions *PermissionSet, validate bool) error {
	selectedProjects, err := selectProjectsFromHarbor()
	if err != nil {
		return fmt.Errorf("error selecting projects: %v", err)
	}
	if len(selectedProjects) == 0 {
		return nil
	}

	fmt.Println("Select permissions to apply to all selected projects:")
	selected, err := prompt.GetRobotPermissionsFromUser("project")
	if err != nil {
		return fmt.Errorf("failed to get permissions: %v", utils.ParseHarborErrorMsg(err))
	}
	if validate {
		if selected, err = ValidateProjectPermissions(selected); err != nil {
			return err
		}
	}

	for _, projectName := range selectedProjects {
		permissions.Project[projectName] = selected
	}
	return nil
}

// collectPerProjectPermissions asks for projects one by one, each with its
// own permission selection.
func collectPerProjectPermissions(permissions *PermissionSet, validate bool) error {
	for {
		projectName, err := prompt.GetProjectNameFromUser()
		if err != nil {
			return fmt.Errorf("%v", utils.ParseHarborErrorMsg(err))
		}
		if projectName == "" {
			return fmt.Errorf("project name cannot be empty")
		}

		selected, err := prompt.GetRobotPermissionsFromUser("project")
		if err != nil {
			return fmt.Errorf("failed to get permissions: %v", utils.ParseHarborErrorMsg(err))
		}
		if validate {
			if selected, err = ValidateProjectPermissions(selected); err != nil {
				return err
			}
		}
		permissions.Project[projectName] = selected

		more, err := robotview.ConfirmMoreProjects()
		if err != nil {
			return fmt.Errorf("error asking for more projects: %v", err)
		}
		if !more {
			return nil
		}
	}
}

// selectProjectsFromHarbor lists all projects and asks which ones to grant
// shared permissions on.
func selectProjectsFromHarbor() ([]string, error) {
	allProjects, err := api.ListAllProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %v", err)
	}

	projectNames := make([]string, 0, len(allProjects.Payload))
	for _, project := range allProjects.Payload {
		projectNames = append(projectNames, project.Name)
	}
	return robotview.SelectProjectsForPermissions(projectNames)
}
