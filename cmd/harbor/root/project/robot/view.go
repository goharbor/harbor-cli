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
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/robot"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/view"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func ViewRobotCommand() *cobra.Command {
	var ProjectName string
	cmd := &cobra.Command{
		Use:   "view [robotID]",
		Short: "get robot by id",
		Long: `View detailed information about a robot account in Harbor.

This command displays comprehensive information about a robot account including
its ID, name, description, creation time, expiration, and the permissions
it has been granted within its project.

The command supports multiple ways to identify the robot account:
- By providing the robot ID directly as an argument
- By specifying a project with the --project flag and selecting the robot interactively
- Without any arguments, which will prompt for both project and robot selection

The displayed information includes:
- Basic details (ID, name, description)
- Temporal information (creation date, expiration date, remaining time)
- Security details (disabled status)
- Detailed permissions breakdown by resource and action

Examples:
  # View robot by ID
  harbor-cli project robot view 123

  # View robot by selecting from a specific project
  harbor-cli project robot view --project myproject

  # Interactive selection (will prompt for project and robot)
  harbor-cli project robot view`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				robot   *robot.GetRobotByIDOK
				robotID int64
				err     error
			)

			if len(args) == 1 {
				robotID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Fatalf("failed to parse robot ID: %v", err)
				}
			} else if ProjectName != "" {
				project, err := api.GetProject(ProjectName, false)
				if err != nil {
					log.Fatalf("failed to get project by name %s: %v", ProjectName, err)
				}
				robotID = prompt.GetRobotIDFromUser(int64(project.Payload.ProjectID))
			} else {
				projectID, err := prompt.GetProjectIDFromUser()
				if err != nil {
					log.Fatalf("failed to get project by id %s: %v", projectID, utils.ParseHarborErrorMsg(err))
				}
				robotID = prompt.GetRobotIDFromUser(projectID)
			}

			robot, err = api.GetRobot(robotID)
			if err != nil {
				log.Fatalf("failed to get robot: %v", err)
			}

			// Convert to a list and display
			// robots := &models.Robot{robot.Payload}
			view.ViewRobot(robot.Payload)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&ProjectName, "project", "", "", "set project name")
	return cmd
}
