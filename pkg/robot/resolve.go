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
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

// ResolveRobotID parses the robot ID argument, or walks the user through
// selecting one. For project-level robots the project comes from the
// --project flag or an interactive selection.
func ResolveRobotID(args []string, level Level, projectName string) (int64, error) {
	if len(args) == 1 {
		robotID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse robot ID: %v", err)
		}
		return robotID, nil
	}
	return promptRobotID(level, projectName)
}

func promptRobotID(level Level, projectName string) (int64, error) {
	projectID := int64(-1)
	if level == LevelProject {
		var err error
		projectID, err = resolveProjectID(projectName)
		if err != nil {
			return 0, err
		}
	}

	robotID, err := prompt.GetRobotIDFromUser(projectID)
	if err != nil {
		return 0, fmt.Errorf("failed to get robot ID from user: %v", utils.ParseHarborErrorMsg(err))
	}
	return robotID, nil
}

func resolveProjectID(projectName string) (int64, error) {
	if projectName == "" {
		projectID, err := prompt.GetProjectIDFromUser()
		if err != nil {
			return 0, fmt.Errorf("failed to get project from user: %v", utils.ParseHarborErrorMsg(err))
		}
		return projectID, nil
	}

	project, err := api.GetProject(projectName, false)
	if err != nil {
		return 0, fmt.Errorf("failed to get project by name %s: %v", projectName, utils.ParseHarborErrorMsg(err))
	}
	return int64(project.Payload.ProjectID), nil
}

// ResolveRobotIDByName looks a robot up by its full name. Project-level
// lookups need a project, which is asked for when not given.
func ResolveRobotIDByName(robotName string, level Level, projectName string) (int64, error) {
	if level == LevelProject && projectName == "" {
		projectID, err := prompt.GetProjectIDFromUser()
		if err != nil {
			return 0, fmt.Errorf("failed to get project from user: %v", utils.ParseHarborErrorMsg(err))
		}
		project, err := api.GetProject(strconv.FormatInt(projectID, 10), true)
		if err != nil {
			return 0, fmt.Errorf("failed to get project by id %d: %v", projectID, utils.ParseHarborErrorMsg(err))
		}
		projectName = project.Payload.Name
	}

	var qualifier []string
	if level == LevelProject {
		qualifier = append(qualifier, projectName)
	}

	robot, err := api.GetRobotByName(robotName, qualifier...)
	if err != nil {
		return 0, WrapError("get robot by name "+robotName, err)
	}
	return robot.ID, nil
}
