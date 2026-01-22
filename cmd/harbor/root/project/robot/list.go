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
		Long: `List robot accounts in a Harbor project.

This command displays a list of robot accounts, either from a specific project
or by prompting you to select a project interactively. The list includes basic
information about each robot account, such as ID, name, creation time, and
expiration status.

The command supports multiple ways to specify the project:
- By providing a project name as an argument
- By using the --project-id flag
- By using the -q/--query flag with a project filter
- Without any arguments, which will prompt for project selection

You can control the output using pagination flags and format options:
- Use --page and --page-size to navigate through results
- Use --sort to order the results
- Set output-format in your configuration for JSON, YAML, or other formats

Examples:
  # List robots in a specific project by name
  harbor-cli project robot list myproject

  # List robots in a project by ID
  harbor-cli project robot list --project-id 123

  # List robots with pagination
  harbor-cli project robot list --page 2 --page-size 20

  # List robots with custom sorting
  harbor-cli project robot list --sort name

  # Interactive listing (will prompt for project selection)
  harbor-cli project robot list`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}
			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

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
				projectID, err := prompt.GetProjectIDFromUser()
				if err != nil {
					log.Fatalf("failed to get project by id %d: %v", projectID, utils.ParseHarborErrorMsg(err))
				}
				opts.Q = projectQString + strconv.FormatInt(projectID, 10)
			}

			robots, err := api.ListRobot(opts)
			if err != nil {
				log.Errorf("failed to get robots list: %v", err)
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				log.WithField("output_format", formatFlag).Debug("Output format selected")
				err = utils.PrintFormat(robots, formatFlag)
				if err != nil {
					return err
				}
			} else {
				list.ListRobots(robots.Payload)
			}
			return nil
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
