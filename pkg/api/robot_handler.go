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
package api

import (
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/permissions"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/robot"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	"github.com/goharbor/harbor-cli/pkg/views/robot/update"
	log "github.com/sirupsen/logrus"
)

func ListRobot(opts ListFlags) (*robot.ListRobotOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Robot.ListRobot(
		ctx,
		&robot.ListRobotParams{
			Page:     &opts.Page,
			PageSize: &opts.PageSize,
			Q:        &opts.Q,
			Sort:     &opts.Sort,
		},
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetRobot(robotID int64) (*robot.GetRobotByIDOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Robot.GetRobotByID(ctx, &robot.GetRobotByIDParams{RobotID: robotID})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetRobotByName(targetRobotName string, projectName ...string) (*robot.GetRobotByIDOK, error) {
	var listResponse *robot.ListRobotOK
	var err error
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	if len(projectName) > 0 {
		var project *project.GetProjectOK
		project, err = GetProject(projectName[0], false)
		if err != nil {
			return nil, fmt.Errorf("failed to get project: %v", utils.ParseHarborErrorMsg(err))
		}
		listResponse, err = ListRobot(ListFlags{
			Q: constants.ProjectQString + strconv.FormatInt(int64(project.Payload.ProjectID), 10),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list robots: %v", utils.ParseHarborErrorMsg(err))
		}
	} else {
		listResponse, err = ListRobot(ListFlags{})
		if err != nil {
			return nil, fmt.Errorf("failed to list robots: %v", utils.ParseHarborErrorMsg(err))
		}
	}

	robotID := int64(-1)
	for _, robotItem := range listResponse.Payload {
		if robotItem.Name == targetRobotName {
			robotID = robotItem.ID
			break
		}
	}
	if robotID == -1 {
		return nil, fmt.Errorf("failed to find robot with name: %v, does it exist?", targetRobotName)
	}

	response, err := client.Robot.GetRobotByID(ctx, &robot.GetRobotByIDParams{RobotID: robotID})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func CheckRoboWithNameExists(projectID int32, name string) (bool, error) {
	var exists bool = false

	response, err := ListRobot(ListFlags{
		Q: constants.ProjectQString + strconv.FormatInt(int64(projectID), 10),
	})
	if err != nil {
		return exists, fmt.Errorf("failed to list robots: %v", utils.ParseHarborErrorMsg(err))
	}

	project, err := GetProject(strconv.FormatInt(int64(projectID), 10), true)
	if err != nil {
		return false, fmt.Errorf("failed to get project: %v", utils.ParseHarborErrorMsg(err))
	}
	projectName := project.Payload.Name

	configurations, err := GetConfigurations()
	if err != nil {
		return exists, fmt.Errorf("failed to get configurations: %v", utils.ParseHarborErrorMsg(err))
	}
	robotNamePrefix := configurations.Payload.RobotNamePrefix.Value
	targetRobotName := fmt.Sprintf("%s%s+%s", robotNamePrefix, projectName, name)
	for _, robot := range response.Payload {
		if robot.Name == targetRobotName {
			exists = true
			break
		}
	}

	return exists, nil
}

func DeleteRobot(robotID int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Robot.DeleteRobot(ctx, &robot.DeleteRobotParams{RobotID: robotID})
	if err != nil {
		return err
	}

	log.Info("robot deleted successfully")
	return nil
}

func CreateRobot(opts create.CreateView) (*robot.CreateRobotCreated, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	// Create a slice to store converted permissions
	permissions := opts.Permissions
	convertedPerms := make([]*models.RobotPermission, 0, len(permissions))

	// Loop through original permissions and convert them
	for _, perm := range permissions {
		convertedPerm := &models.RobotPermission{
			Access:    perm.Access,
			Kind:      perm.Kind,
			Namespace: perm.Namespace,
		}
		convertedPerms = append(convertedPerms, convertedPerm)
	}

	response, err := client.Robot.CreateRobot(
		ctx,
		&robot.CreateRobotParams{
			Robot: &models.RobotCreate{
				Description: opts.Description,
				Disable:     false,
				Duration:    opts.Duration,
				Level:       opts.Level,
				Name:        opts.Name,
				Permissions: convertedPerms,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create robot: %v", utils.ParseHarborErrorMsg(err))
	}

	log.Info("robot created successfully.")
	return response, nil
}

// update robot with robotID
func UpdateRobot(opts *update.UpdateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		log.Errorf("Error: %v", err)
		return err
	}

	// Create a slice to store converted permissions
	permissions := opts.Permissions
	convertedPerms := make([]*models.RobotPermission, 0, len(permissions))

	for _, perm := range permissions {
		convertedPerm := &models.RobotPermission{
			Access:    perm.Access,
			Kind:      perm.Kind,
			Namespace: perm.Namespace,
		}
		convertedPerms = append(convertedPerms, convertedPerm)
	}

	_, err = client.Robot.UpdateRobot(
		ctx,
		&robot.UpdateRobotParams{
			Robot: &models.Robot{
				Description: opts.Description,
				Duration:    &opts.Duration,
				Editable:    opts.Editable,
				Disable:     opts.Disable,
				ID:          opts.ID,
				Level:       opts.Level,
				Name:        opts.Name,
				Permissions: convertedPerms,
			},
			RobotID: opts.ID,
		},
	)
	if err != nil {
		log.Errorf("Error in updating Robot: %v", err)
		return err
	}

	log.Info("robot updated successfully.")
	return nil
}

func GetPermissions() (*permissions.GetPermissionsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Permissions.GetPermissions(
		ctx,
		&permissions.GetPermissionsParams{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func RefreshSecret(secret string, robotID int64) (*robot.RefreshSecOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	robotSec := &models.RobotSec{
		Secret: secret,
	}

	response, err := client.Robot.RefreshSec(ctx, &robot.RefreshSecParams{
		RobotSec: robotSec,
		RobotID:  robotID,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
