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
	baseprompt "github.com/goharbor/harbor-cli/pkg/prompt"
	robotprompt "github.com/goharbor/harbor-cli/pkg/prompt/robot"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	"github.com/sirupsen/logrus"
)

// getProjectPermissions orchestrates project permission selection for create or update flows.
func getProjectPermissions(isUpdate bool, createOpts *create.CreateView, projectPermissionsMap map[string][]models.Permission) error {
	if isUpdate {
		return handleProjectPermissionsUpdate(projectPermissionsMap)
	}
	return handleProjectPermissionsCreate(createOpts, projectPermissionsMap)
}

// handleProjectPermissionsCreate handles the create mode flow.
func handleProjectPermissionsCreate(createOpts *create.CreateView, projectPermissionsMap map[string][]models.Permission) error {
	if createOpts == nil {
		return fmt.Errorf("create options must be provided for create flow")
	}

	permissionMode, err := robotprompt.ChooseProjectPermissionMode(false)
	if err != nil {
		return fmt.Errorf("error selecting permission mode: %v", err)
	}

	switch permissionMode {
	case "list":
		return assignCommonPermissionsToSelectedProjects(projectPermissionsMap)
	case "per_project":
		return assignPerProjectPermissionsInteractive(createOpts, projectPermissionsMap)
	case "none", "clear":
		fmt.Println("Creating robot with system-level permissions only (no project-specific permissions)")
		return nil
	default:
		return fmt.Errorf("unknown permission mode: %s", permissionMode)
	}
}

// handleProjectPermissionsUpdate handles the update mode flow.
func handleProjectPermissionsUpdate(projectPermissionsMap map[string][]models.Permission) error {
	hasExisting := len(projectPermissionsMap) > 0
	permissionMode, err := robotprompt.ChooseProjectPermissionMode(hasExisting)
	if err != nil {
		return fmt.Errorf("error selecting permission mode: %v", err)
	}

	switch permissionMode {
	case "keep":
		logrus.Info("Keeping existing project permissions")
		return nil
	case "clear":
		logrus.Info("Clearing all project permissions")
		clearProjectPermissions(projectPermissionsMap)
		return nil
	case "list":
		return assignCommonPermissionsToSelectedProjectsForUpdate(projectPermissionsMap)
	case "per_project":
		return modifyPerProjectPermissionsInteractive(projectPermissionsMap)
	default:
		return fmt.Errorf("unknown permission mode: %s", permissionMode)
	}
}

// assignCommonPermissionsToSelectedProjects prompts for projects, then one permission set to apply to all.
func assignCommonPermissionsToSelectedProjects(projectPermissionsMap map[string][]models.Permission) error {
	selectedProjects, err := baseprompt.GetProjectNamesFromUser()
	if err != nil {
		return fmt.Errorf("error selecting projects: %v", err)
	}
	if len(selectedProjects) == 0 {
		return nil
	}
	perms := promptCommonProjectPermissions()
	setPermissionsForProjects(selectedProjects, perms, projectPermissionsMap)
	return nil
}

// assignPerProjectPermissionsInteractive lets the user add one or many projects, each with its own permissions.
func assignPerProjectPermissionsInteractive(opts *create.CreateView, projectPermissionsMap map[string][]models.Permission) error {
	if opts.ProjectName != "" {
		projectPermissionsMap[opts.ProjectName] = baseprompt.GetRobotPermissionsFromUser("project")
		return nil
	}
	return addProjectsInteractively(projectPermissionsMap)
}

// assignCommonPermissionsToSelectedProjectsForUpdate handles update mode, optionally replacing existing first.
func assignCommonPermissionsToSelectedProjectsForUpdate(projectPermissionsMap map[string][]models.Permission) error {
	if len(projectPermissionsMap) > 0 {
		replaceExisting, err := robotprompt.ConfirmReplaceExisting()
		if err != nil {
			return fmt.Errorf("error asking about existing permissions: %v", err)
		}
		if replaceExisting {
			clearProjectPermissions(projectPermissionsMap)
		}
	}
	return assignCommonPermissionsToSelectedProjects(projectPermissionsMap)
}

func modifyPerProjectPermissionsInteractive(projectPermissionsMap map[string][]models.Permission) error {
	if len(projectPermissionsMap) > 0 {
		modifyMode, err := robotprompt.ChooseModifyMode()
		if err != nil {
			return fmt.Errorf("error asking about permission modification: %v", err)
		}
		switch modifyMode {
		case "replace":
			clearProjectPermissions(projectPermissionsMap)
		case "modify":
			selected, err := selectExistingProjectsForModify(projectPermissionsMap)
			if err != nil {
				return err
			}
			for _, p := range selected {
				fmt.Printf("Updating permissions for project: %s\n", p)
				projectPermissionsMap[p] = baseprompt.GetRobotPermissionsFromUser("project")
			}
			return nil
		case "add":
			// fall through to add loop below
		default:
			return fmt.Errorf("unknown modify mode: %s", modifyMode)
		}
	}
	return addProjectsInteractively(projectPermissionsMap)
}

func promptCommonProjectPermissions() []models.Permission {
	return baseprompt.GetRobotPermissionsFromUser("project")
}

func setPermissionsForProjects(projects []string, perms []models.Permission, projectPermissionsMap map[string][]models.Permission) {
	for _, name := range projects {
		projectPermissionsMap[name] = perms
	}
}

func clearProjectPermissions(projectPermissionsMap map[string][]models.Permission) {
	for k := range projectPermissionsMap {
		delete(projectPermissionsMap, k)
	}
}

func selectExistingProjectsForModify(projectPermissionsMap map[string][]models.Permission) ([]string, error) {
	provider := func() ([]string, error) {
		list := make([]string, 0, len(projectPermissionsMap))
		for p := range projectPermissionsMap {
			list = append(list, p)
		}
		return list, nil
	}
	return robotprompt.SelectProjects(provider)
}

func addProjectsInteractively(projectPermissionsMap map[string][]models.Permission) error {
	for {
		projectName, err := baseprompt.GetProjectNameFromUser()
		if err != nil {
			return fmt.Errorf("%v", utils.ParseHarborErrorMsg(err))
		}
		if projectName == "" {
			return fmt.Errorf("project name cannot be empty")
		}
		projectPermissionsMap[projectName] = baseprompt.GetRobotPermissionsFromUser("project")
		moreProjects, err := robotprompt.AskMoreProjects()
		if err != nil {
			return fmt.Errorf("error asking for more projects: %v", err)
		}
		if !moreProjects {
			return nil
		}
	}
}
