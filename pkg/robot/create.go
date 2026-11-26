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
	"github.com/spf13/viper"
)

// CreateOptions parameterizes RunCreate for both robot levels.
type CreateOptions struct {
	Level        Level
	View         *robotview.CreateView
	All          bool
	ConfigFile   string
	ExportToFile bool
}

// RunCreate creates a robot account from a config file, flags, or
// interactive input.
func RunCreate(opts CreateOptions) error {
	view := opts.View

	permissions, err := gatherCreateInput(opts)
	if err != nil {
		return err
	}

	view.Permissions = permissions.Merge()
	view.Level = opts.Level.String()
	if len(view.Permissions) == 0 {
		return errNoPermissions("create")
	}

	if err := ensureNameAvailable(opts.Level, view.ProjectName, view.Name); err != nil {
		return err
	}

	response, err := api.CreateRobot(*view)
	if err != nil {
		return WrapError("create robot", err)
	}
	return showCreatedRobot(response.Payload, opts.ExportToFile)
}

func gatherCreateInput(opts CreateOptions) (*PermissionSet, error) {
	if opts.ConfigFile != "" {
		loaded, permissions, err := LoadCreateFromConfig(opts.Level, opts.ConfigFile)
		if err != nil {
			return nil, err
		}
		*opts.View = *loaded
		return permissions, nil
	}

	if opts.Level == LevelProject {
		if err := ensureProjectName(opts.View); err != nil {
			return nil, err
		}
	}
	if err := promptBasicInfo(opts.View); err != nil {
		return nil, err
	}
	return collectCreatePermissions(opts.Level, opts.View, opts.All)
}

func ensureProjectName(view *robotview.CreateView) error {
	if view.ProjectName != "" {
		return nil
	}

	projectName, err := prompt.GetProjectNameFromUser()
	if err != nil {
		return fmt.Errorf("%v", utils.ParseHarborErrorMsg(err))
	}
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	view.ProjectName = projectName
	return nil
}

// ensureNameAvailable rejects names that already exist at the given level,
// so the user gets a clear error instead of Harbor's conflict response.
func ensureNameAvailable(level Level, projectName, name string) error {
	var exists bool
	var err error

	if level == LevelProject {
		project, projectErr := api.GetProject(projectName, false)
		if projectErr != nil {
			return fmt.Errorf("failed to get project: %v", utils.ParseHarborErrorMsg(projectErr))
		}
		exists, err = api.CheckRobotWithNameExists(project.Payload.ProjectID, name)
	} else {
		exists, err = api.CheckSystemRobotWithNameExists(name)
	}

	if err != nil {
		return fmt.Errorf("failed to get robot by name: %v", utils.ParseHarborErrorMsg(err))
	}
	if exists {
		if level == LevelProject {
			return fmt.Errorf("robot account with name '%s' already exists in project '%s'", name, projectName)
		}
		return fmt.Errorf("robot account with name '%s' already exists", name)
	}
	return nil
}

func showCreatedRobot(robot *models.RobotCreated, exportToFile bool) error {
	fmt.Printf("Successfully created robot account '%s' (ID: %d)\n", robot.Name, robot.ID)

	if viper.GetString("output-format") != "" {
		details, err := api.GetRobot(robot.ID)
		if err != nil {
			return fmt.Errorf("failed to get robot details: %v", utils.ParseHarborErrorMsg(err))
		}
		utils.SavePayloadJSON(robot.Name, details.Payload)
		return nil
	}

	if exportToFile {
		fmt.Printf("Exporting robot credentials to file\n")
		return ExportSecretToFile(robot.Name, robot.Secret, robot.CreationTime.String(), robot.ExpiresAt)
	}
	return DisplaySecret(robot.Name, robot.Secret)
}
