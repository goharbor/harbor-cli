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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListRobotCommand creates a new `harbor project robot list` command
func ListRobotCommand() *cobra.Command {
	var opts api.ListFlags

	projectQString := constants.ProjectQString
	cmd := &cobra.Command{
		Use:   "list [projectName]",
		Short: "list robot",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				project, err := api.GetProject(args[0], false)
				if err != nil {
					log.Errorf("Invalid Project Name: %v", err)
				}
				opts.ProjectID = int64(project.Payload.ProjectID)
				opts.Q = projectQString + strconv.FormatInt(opts.ProjectID, 10)
			} else if opts.Q != "" {
				opts.Q = projectQString + opts.Q
			} else if opts.ProjectID > 0 {
				opts.Q = projectQString + strconv.FormatInt(opts.ProjectID, 10)
			} else {
				projectID := prompt.GetProjectIDFromUser()
				opts.Q = projectQString + strconv.FormatInt(projectID, 10)
			}

			robots, err := api.ListRobot(opts)
			if err != nil {
				log.Errorf("failed to get robots list: %v", err)
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(robots, formatFlag)
				if err != nil {
					log.Errorf("Invalid Print Format: %v", err)
				}
			} else {
				list.ListRobots(robots.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.Int64VarP(&opts.ProjectID, "project-id", "", 0, "Project ID")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(
		&opts.Sort,
		"sort",
		"",
		"",
		"Sort the resource list in ascending or descending order",
	)

	return cmd
}
