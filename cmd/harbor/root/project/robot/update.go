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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/update"

	"github.com/spf13/cobra"
)

func UpdateRobotCommand() *cobra.Command {
	var (
		robotID     int64
		opts        update.UpdateView
		all         bool
		ProjectName string
	)

	cmd := &cobra.Command{
		Use:   "update [robotID]",
		Short: "update robot by id",
		Long: `Update an existing robot account within a Harbor project.

Robot accounts are non-human users that can be used for automation purposes
such as CI/CD pipelines, scripts, or other automated processes that need
to interact with Harbor. This command allows you to modify an existing robot's
properties including its name, description, duration, and permissions.

This command supports both interactive and non-interactive modes:
- With robot ID: directly updates the specified robot
- With --project flag: helps select a robot from the specified project
- Without either: walks through project and robot selection interactively

The update process will:
1. Identify the robot account to be updated
2. Load its current configuration
3. Apply the requested changes
4. Save the updated configuration

Fields that can be updated:
- Name: The robot account's identifier
- Description: A human-readable description of the robot's purpose
- Duration: The lifetime of the robot account in days
- Permissions: The actions the robot is allowed to perform

Note: Updating a robot does not regenerate its secret. If you need a new
secret, consider deleting the robot and creating a new one instead.

Examples:
  # Update robot by ID with a new description
  harbor-cli project robot update 123 --description "Updated CI/CD pipeline robot"

  # Update robot's duration (extend lifetime)
  harbor-cli project robot update 123 --duration 180

  # Update by selecting from a specific project
  harbor-cli project robot update --project myproject

  # Update with all permissions
  harbor-cli project robot update 123 --all-permission

  # Interactive update (will prompt for robot selection and changes)
  harbor-cli project robot update`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) == 1 {
				robotID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse robot ID: %v", err)
				}
			} else if ProjectName != "" {
				project, err := api.GetProject(ProjectName, false)
				if err != nil {
					return fmt.Errorf("failed to get project by name %s: %v", ProjectName, err)
				}
				robotID = prompt.GetRobotIDFromUser(int64(project.Payload.ProjectID))
			} else {
				projectID, err := prompt.GetProjectIDFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project by id %d: %v", projectID, utils.ParseHarborErrorMsg(err))
				}
				robotID = prompt.GetRobotIDFromUser(projectID)
			}

			robot, err := api.GetRobot(robotID)
			if err != nil {
				return fmt.Errorf("failed to get robot: %v", err)
			}

			bot := robot.Payload

			var duration int64
			if bot.Duration != nil {
				duration = *bot.Duration
			}

			opts = update.UpdateView{
				CreationTime: bot.CreationTime,
				Description:  bot.Description,
				Disable:      bot.Disable,
				Duration:     duration,
				Editable:     bot.Editable,
				ID:           bot.ID,
				Level:        bot.Level,
				Name:         bot.Name,
				Secret:       bot.Secret,
			}

			// declare empty permissions to hold permissions
			var permissions []models.Permission

			if all {
				perms, _ := api.GetPermissions()
				permission := perms.Payload.Project

				choices := []models.Permission{}
				for _, perm := range permission {
					choices = append(choices, *perm)
				}
				permissions = choices
			} else {
				permissions, err = prompt.GetRobotPermissionsFromUser("project")
				if err != nil {
					return fmt.Errorf("failed to get permissions: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			// []Permission to []*Access
			var accesses []*models.Access
			for _, perm := range permissions {
				access := &models.Access{
					Action:   perm.Action,
					Resource: perm.Resource,
				}
				accesses = append(accesses, access)
			}
			// convert []models.permission to []*model.Access
			perm := &update.RobotPermission{
				Kind:      bot.Permissions[0].Kind,
				Namespace: bot.Permissions[0].Namespace,
				Access:    accesses,
			}
			opts.Permissions = []*update.RobotPermission{perm}

			err = updateRobotView(&opts)
			if err != nil {
				return fmt.Errorf("failed to Update robot: %v", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(
		&all,
		"all-permission",
		"a",
		false,
		"Select all permissions for the robot account",
	)
	flags.StringVarP(&opts.Name, "name", "", "", "name of the robot account")
	flags.StringVarP(&opts.Description, "description", "", "", "description of the robot account")
	flags.StringVarP(&ProjectName, "project", "", "", "set project name")
	flags.Int64VarP(&opts.Duration, "duration", "", 0, "set expiration of robot account in days")

	return cmd
}

func updateRobotView(updateView *update.UpdateView) error {
	if updateView == nil {
		updateView = &update.UpdateView{}
	}

	update.UpdateRobotView(updateView)
	return api.UpdateRobot(updateView)
}
