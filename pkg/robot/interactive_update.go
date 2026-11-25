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

	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/sirupsen/logrus"
)

// promptUpdate runs the interactive update flow: the basic form, then an
// optional pass over the robot's permissions.
func promptUpdate(view *robotview.UpdateView, level Level, all bool, permissions *PermissionSet) error {
	if err := robotview.FormUpdate(view); err != nil {
		return fmt.Errorf("failed to read robot input: %v", err)
	}
	if view.Duration == 0 {
		return fmt.Errorf("failed to update robot: duration cannot be 0")
	}

	updatePermissions, err := robotview.ConfirmUpdatePermissions()
	if err != nil {
		return fmt.Errorf("error asking about permission updates: %v", err)
	}
	if !updatePermissions {
		logrus.Debug("Keeping existing permissions")
		return nil
	}

	if level == LevelProject {
		return updateSingleProjectPermissions(all, permissions)
	}

	if err := updateSystemPermissions(all, permissions); err != nil {
		return err
	}
	return updateProjectPermissions(permissions)
}

// updateSingleProjectPermissions replaces the one scope of a project-level
// robot.
func updateSingleProjectPermissions(all bool, permissions *PermissionSet) error {
	selected, err := collectSingleProjectPermissions(all, "update")
	if err != nil {
		return err
	}
	for projectName := range permissions.Project {
		permissions.Project[projectName] = selected
	}
	return nil
}

func updateSystemPermissions(all bool, permissions *PermissionSet) error {
	updateSystem, err := robotview.ConfirmUpdateSystemPermissions()
	if err != nil {
		return fmt.Errorf("error asking about system permission updates: %v", err)
	}
	if !updateSystem {
		logrus.Debug("Keeping existing system permissions")
		return nil
	}

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
		return fmt.Errorf("failed to update robot: %v", utils.ParseHarborErrorMsg(err))
	}
	if len(selected) == 0 {
		return errNoPermissions("update")
	}
	permissions.System = selected
	return nil
}

func updateProjectPermissions(permissions *PermissionSet) error {
	mode, err := robotview.SelectUpdatePermissionMode(len(permissions.Project) > 0)
	if err != nil {
		return fmt.Errorf("error selecting permission mode: %v", err)
	}

	switch mode {
	case robotview.ModeKeep:
		logrus.Debug("Keeping existing project permissions")
		return nil
	case robotview.ModeClear:
		logrus.Debug("Clearing all project permissions")
		permissions.ClearProjects()
		return nil
	case robotview.ModeList:
		return updateMultiProjectPermissions(permissions)
	case robotview.ModePerProject:
		return updatePerProjectPermissions(permissions)
	default:
		return fmt.Errorf("unknown permission mode: %s", mode)
	}
}

func updateMultiProjectPermissions(permissions *PermissionSet) error {
	if len(permissions.Project) > 0 {
		replace, err := robotview.ConfirmReplaceExistingProjects()
		if err != nil {
			return fmt.Errorf("error asking about existing permissions: %v", err)
		}
		if replace {
			permissions.ClearProjects()
		}
	}
	return collectMultiProjectPermissions(permissions, true)
}

func updatePerProjectPermissions(permissions *PermissionSet) error {
	if len(permissions.Project) == 0 {
		return collectPerProjectPermissions(permissions, true)
	}

	mode, err := robotview.SelectModifyMode()
	if err != nil {
		return fmt.Errorf("error asking about permission modification: %v", err)
	}

	switch mode {
	case robotview.ModifyReplace:
		permissions.ClearProjects()
		return collectPerProjectPermissions(permissions, true)
	case robotview.ModifyModify:
		return modifyExistingProjects(permissions)
	case robotview.ModifyAdd:
		return collectPerProjectPermissions(permissions, true)
	default:
		return fmt.Errorf("unknown permission mode: %s", mode)
	}
}

// modifyExistingProjects re-permissions a selection of the projects the
// robot already has access to.
func modifyExistingProjects(permissions *PermissionSet) error {
	selectedProjects, err := robotview.SelectProjectsToModify(permissions.ProjectNames())
	if err != nil {
		return fmt.Errorf("error selecting projects to modify: %v", err)
	}

	validate, err := projectPermissionValidator()
	if err != nil {
		return err
	}

	for _, projectName := range selectedProjects {
		fmt.Printf("Updating permissions for project: %s\n", projectName)
		selected, err := prompt.GetRobotPermissionsFromUser("project")
		if err != nil {
			return fmt.Errorf("failed to get permissions: %v", utils.ParseHarborErrorMsg(err))
		}
		valid := validate(selected)
		if len(valid) == 0 {
			return errNoPermissions("update")
		}
		permissions.Project[projectName] = valid
	}
	return nil
}
