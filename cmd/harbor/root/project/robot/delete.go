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

	"github.com/goharbor/harbor-cli/pkg/api"
	robot "github.com/goharbor/harbor-cli/pkg/robot"
	"github.com/spf13/cobra"
)

// to-do improve DeleteRobotCommand and multi select & delete
func DeleteRobotCommand() *cobra.Command {
	var projectName string
	cmd := &cobra.Command{
		Use:   "delete [robotName]",
		Short: "delete robot by name",
		Long: `Delete a robot account from a Harbor project.

This command permanently removes a robot account from Harbor. Once deleted,
the robot's credentials will no longer be valid, and any automated processes
using those credentials will fail.

The command supports multiple ways to identify the robot account to delete:
- By providing the robot ID directly as an argument
- By specifying a project with the --project flag and selecting the robot interactively
- Without any arguments, which will prompt for both project and robot selection

Important considerations:
- Deletion is permanent and cannot be undone
- All access tokens for the robot will be invalidated immediately
- Any systems using the robot's credentials will need to be updated

Examples:
  # Delete robot by Name, choose project
  harbor-cli project robot delete robot_projectname+robotname

  # Delete robot by Name and project name
  harbor-cli project robot delete robot_projectname+robotname --project projectname

  # Delete robot by selecting from a specific project
  harbor-cli project robot delete --project myproject

  # Interactive deletion (will prompt for project and robot selection)
  harbor-cli project robot delete`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				robotID int64
				err     error
			)
			if len(args) == 1 {
				robotID, err = robot.ResolveRobotIDByName(args[0], robot.LevelProject, projectName)
			} else {
				robotID, err = robot.ResolveRobotID(args, robot.LevelProject, projectName)
			}
			if err != nil {
				return err
			}

			if err := api.DeleteRobot(robotID); err != nil {
				return robot.WrapError("delete robot", err)
			}
			fmt.Printf("Robot account (ID: %d) was successfully deleted\n", robotID)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&projectName, "project", "", "", "set project name")
	return cmd
}
