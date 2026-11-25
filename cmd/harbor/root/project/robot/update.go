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
	robot "github.com/goharbor/harbor-cli/pkg/robot"
	"github.com/spf13/cobra"
)

func UpdateRobotCommand() *cobra.Command {
	var (
		all         bool
		configFile  string
		projectName string
		// Bound to flags for CLI compatibility; the interactive flow loads
		// the current values from Harbor instead (see goharbor/harbor-cli#1150).
		name        string
		description string
		duration    int64
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

  # Update from configuration file
  harbor-cli project robot update 123 --robot-config-file ./robot-config.yaml

  # Interactive update (will prompt for robot selection and changes)
  harbor-cli project robot update`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			robotID, err := robot.ResolveRobotID(args, robot.LevelProject, projectName)
			if err != nil {
				return err
			}
			return robot.RunUpdate(robot.UpdateOptions{
				Level:      robot.LevelProject,
				RobotID:    robotID,
				All:        all,
				ConfigFile: configFile,
			})
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&all, "all-permission", "a", false, "Select all permissions for the robot account")
	flags.StringVarP(&name, "name", "", "", "name of the robot account")
	flags.StringVarP(&description, "description", "", "", "description of the robot account")
	flags.StringVarP(&projectName, "project", "", "", "set project name")
	flags.Int64VarP(&duration, "duration", "", 0, "set expiration of robot account in days")
	flags.StringVarP(&configFile, "robot-config-file", "r", "", "YAML/JSON file with robot configuration")

	return cmd
}
