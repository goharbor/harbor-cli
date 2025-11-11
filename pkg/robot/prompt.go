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
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	"github.com/sirupsen/logrus"
)

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

func promptCommonProjectPermissions() []models.Permission {
	return baseprompt.GetRobotPermissionsFromUser("project")
}
