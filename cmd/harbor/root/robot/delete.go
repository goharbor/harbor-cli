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
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views"

	"github.com/spf13/cobra"
)

// to-do improve DeleteRobotCommand and multi select & delete
func DeleteRobotCommand() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete [robotName]",
		Short: "delete robot by name",
		Long: `Delete a robot account from Harbor.

This command permanently removes a robot account from Harbor. Once deleted,
the robot's credentials will no longer be valid, and any automated processes
using those credentials will fail.

The command supports multiple ways to identify the robot account to delete:
- By providing the robot name directly as an argument
- Without any arguments, which will prompt for robot selection

Important considerations:
- Deletion is permanent and cannot be undone
- All access tokens for the robot will be invalidated immediately
- Any systems using the robot's credentials will need to be updated
- For system robots, access across all projects will be revoked

Examples:
  # Delete robot by name
  harbor-cli robot delete robot_robotname

  # Interactive deletion (will prompt for robot selection)
  harbor-cli robot delete`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				robotID int64
				err     error
				name    string
			)
			if len(args) == 1 {
				robotName := args[0]
				robot, err := api.GetRobotByName(robotName)
				if err != nil {
					errorCode := utils.ParseHarborErrorCode(err)
					if errorCode == "403" {
						return fmt.Errorf("Permission denied: (Project) Admin privileges are required to execute this command.")
					} else {
						return fmt.Errorf("failed to get robot: %v", utils.ParseHarborErrorMsg(err))
					}
				}
				robotID = robot.ID
				name = robot.Name
			} else {
				robotID, err = prompt.GetRobotIDFromUser(-1)
				if err != nil {
					return fmt.Errorf("failed to get robot ID from user: %v", utils.ParseHarborErrorMsg(err))
				}
				name = fmt.Sprintf("ID %d", robotID)
			}

			if !yes {
				confirm, err := views.ConfirmDeletion(fmt.Sprintf("Are you sure you want to delete robot account '%s'?", name))
				if err != nil {
					return err
				}
				if !confirm {
					fmt.Println("Deletion cancelled")
					return nil
				}
			}

			err = api.DeleteRobot(robotID)
			if err != nil {
				errorCode := utils.ParseHarborErrorCode(err)
				if errorCode == "403" {
					return fmt.Errorf("Permission denied: (Project) Admin privileges are required to execute this command.")
				} else {
					return fmt.Errorf("failed to delete robot: %v", utils.ParseHarborErrorMsg(err))
				}
			}
			fmt.Printf("Robot account (ID: %d) was successfully deleted\n", robotID)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&yes, "yes", "y", false, "Answer yes to all questions and do not prompt")

	return cmd
}
