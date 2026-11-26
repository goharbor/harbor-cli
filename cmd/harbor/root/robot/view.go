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
	"github.com/goharbor/harbor-cli/pkg/api"
	robot "github.com/goharbor/harbor-cli/pkg/robot"
	"github.com/goharbor/harbor-cli/pkg/views/robot/view"

	"github.com/spf13/cobra"
)

func ViewRobotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [robotID]",
		Short: "get robot by id",
		Long: `View detailed information about a robot account in Harbor.

This command displays comprehensive information about a robot account including
its ID, name, description, creation time, expiration, and the permissions
it has been granted. Supports both system-level and project-level robot accounts.

The command supports multiple ways to identify the robot account:
- By providing the robot ID directly as an argument
- Without any arguments, which will prompt for robot selection

The displayed information includes:
- Basic details (ID, name, description)
- Temporal information (creation date, expiration date, remaining time)
- Security details (disabled status)
- Detailed permissions breakdown by resource and action
- For system robots: permissions across multiple projects are shown separately

System-level robots can have permissions spanning multiple projects, while
project-level robots are scoped to a single project.

Examples:
  # View robot by ID
  harbor-cli robot view 123

  # Interactive selection (will prompt for robot)
  harbor-cli robot view`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			robotID, err := robot.ResolveRobotID(args, robot.LevelSystem, "")
			if err != nil {
				return err
			}

			response, err := api.GetRobot(robotID)
			if err != nil {
				return robot.WrapError("get robot", err)
			}

			view.ViewRobot(response.Payload)
			return nil
		},
	}

	return cmd
}
