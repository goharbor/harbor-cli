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
		all        bool
		configFile string
		// Bound to flags for CLI compatibility; the interactive flow loads
		// the current values from Harbor instead (see goharbor/harbor-cli#1150).
		name        string
		description string
		duration    int64
	)

	cmd := &cobra.Command{
		Use:   "update [robotID]",
		Short: "update robot by id",
		Long: `Update an existing robot account within Harbor.

Robot accounts are non-human users that can be used for automation purposes
such as CI/CD pipelines, scripts, or other automated processes that need
to interact with Harbor. This command allows you to modify an existing robot's
properties including its name, description, duration, and permissions.

This command supports both interactive and non-interactive modes:
- With robot ID: directly updates the specified robot
- Without ID: walks through robot selection interactively

The update process will:
1. Identify the robot account to be updated
2. Load its current configuration
3. Apply the requested changes
4. Save the updated configuration

This command can update both system and project-specific permissions:
- System permissions apply across the entire Harbor instance
- Project permissions apply to specific projects

Configuration can be loaded from:
- Interactive prompts (default)
- Command line flags
- YAML/JSON configuration file

Note: Updating a robot does not regenerate its secret. If you need a new
secret, consider deleting the robot and creating a new one instead.

Examples:
  # Update robot by ID with a new description
  harbor-cli robot update 123 --description "Updated CI/CD pipeline robot"

  # Update robot's duration (extend lifetime)
  harbor-cli robot update 123 --duration 180

  # Update with all permissions
  harbor-cli robot update 123 --all-permission

  # Update from configuration file
  harbor-cli robot update 123 --robot-config-file ./robot-config.yaml

  # Interactive update (will prompt for robot selection and changes)
  harbor-cli robot update`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			robotID, err := robot.ResolveRobotID(args, robot.LevelSystem, "")
			if err != nil {
				return err
			}
			return robot.RunUpdate(robot.UpdateOptions{
				Level:      robot.LevelSystem,
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
	flags.Int64VarP(&duration, "duration", "", 0, "set expiration of robot account in days")
	flags.StringVarP(&configFile, "robot-config-file", "r", "", "YAML/JSON file with robot configuration")

	return cmd
}
