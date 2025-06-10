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
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// to-do improve DeleteRobotCommand and multi select & delete
func DeleteRobotCommand() *cobra.Command {
	var ProjectName string
	cmd := &cobra.Command{
		Use:   "delete [robotID]",
		Short: "delete robot by id",
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
  # Delete robot by ID
  harbor-cli project robot delete 123

  # Delete robot by selecting from a specific project
  harbor-cli project robot delete --project myproject

  # Interactive deletion (will prompt for project and robot selection)
  harbor-cli project robot delete`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				robotID int64
				err     error
			)
			if len(args) == 1 {
				robotID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Fatalf("failed to parse robot ID: %v", utils.ParseHarborErrorMsg(err))
				}
			} else if ProjectName != "" {
				project, err := api.GetProject(ProjectName, false)
				if err != nil {
					log.Fatalf("failed to get project by name %s: %v", ProjectName, utils.ParseHarborErrorMsg(err))
				}
				robotID = prompt.GetRobotIDFromUser(int64(project.Payload.ProjectID))
			} else {
				projectID := prompt.GetProjectIDFromUser()
				robotID = prompt.GetRobotIDFromUser(projectID)
			}
			err = api.DeleteRobot(robotID)
			if err != nil {
				fmt.Printf("failed to delete robots: %v", utils.ParseHarborErrorMsg(err))
				return
			}
			log.Infof("Successfully deleted robot with ID: %d", robotID)
			fmt.Printf("Robot account (ID: %d) was successfully deleted\n", robotID)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&ProjectName, "project", "", "", "set project name")
	return cmd
}
