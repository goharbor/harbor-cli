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

	config "github.com/goharbor/harbor-cli/pkg/config/robot"
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/sirupsen/logrus"
)

// LoadCreateFromConfig reads a robot spec file for `create` and splits its
// permissions. The file's level must match the command's level.
func LoadCreateFromConfig(level Level, path string) (*robotview.CreateView, *PermissionSet, error) {
	fmt.Println("Loading configuration from: ", path)

	view, err := config.LoadConfigFromFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load robot config from file: %v", err)
	}
	logrus.Debug("Successfully loaded robot configuration")

	if err := checkLevel(level, view.Level, "create"); err != nil {
		return nil, nil, err
	}

	if level == LevelProject && view.ProjectName == "" && len(view.Permissions) > 0 {
		view.ProjectName = view.Permissions[0].Namespace
	}

	permissions := NewPermissionSet()
	permissions.Split(view.Permissions)
	return view, permissions, nil
}

// ApplyConfigToUpdate reads a robot spec file for `update`. The file's
// permissions replace the robot's current ones, so the file describes the
// desired state, exactly as it does for create. Description and duration
// are only overwritten when set; name and level cannot be changed through
// the Harbor API and are left alone.
func ApplyConfigToUpdate(view *robotview.UpdateView, level Level, path string) (*PermissionSet, error) {
	fmt.Println("Loading configuration from: ", path)

	loaded, err := config.LoadConfigFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load robot config from file: %v", err)
	}
	logrus.Debug("Successfully loaded robot configuration")

	if err := checkLevel(level, loaded.Level, "update"); err != nil {
		return nil, err
	}

	if loaded.Description != "" {
		view.Description = loaded.Description
	}
	if loaded.Duration != 0 {
		view.Duration = loaded.Duration
	}

	permissions := NewPermissionSet()
	permissions.Split(loaded.Permissions)
	for projectName, projectPermissions := range permissions.Project {
		valid, err := ValidateProjectPermissions(projectPermissions)
		if err != nil {
			return nil, err
		}
		permissions.Project[projectName] = valid
	}
	return permissions, nil
}

func checkLevel(want Level, got, action string) error {
	if got == want.String() {
		return nil
	}
	other := LevelProject
	if want == LevelProject {
		other = LevelSystem
	}
	command := fmt.Sprintf("harbor-cli robot %s", action)
	if other == LevelProject {
		command = fmt.Sprintf("harbor-cli project robot %s", action)
	}
	return fmt.Errorf("invalid robot configuration: level must be '%s'. If you try to %s a %s-level robot, please run the `%s` command instead",
		want, action, other, command)
}
