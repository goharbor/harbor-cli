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
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListRobotCommand creates a new `harbor project robot list` command
func ListRobotCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "list [projectName]",
		Short: "list robot",
		Long: `List robot accounts in Harbor.

This command displays a list of system-level robot accounts. The list includes basic
information about each robot account, such as ID, name, creation time, and
expiration status.

System-level robots have permissions that can span across multiple projects, making
them suitable for CI/CD pipelines and automation tasks that require access to 
multiple projects in Harbor.

You can control the output using pagination flags and format options:
- Use --page and --page-size to navigate through results
- Use --sort to order the results by name, creation time, etc.
- Use -q/--query to filter robots by specific criteria
- Set output-format in your configuration for JSON, YAML, or other formats

Examples:
  # List all system robots
  harbor-cli robot list

  # List system robots with pagination
  harbor-cli robot list --page 2 --page-size 20

  # List system robots with custom sorting
  harbor-cli robot list --sort name

  # Filter system robots by name
  harbor-cli robot list -q name=ci-robot

  # Get robot details in JSON format
  harbor-cli robot list --output-format json`,
		Args: cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}
			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

			robots, err := api.ListRobot(opts)
			if err != nil {
				log.Errorf("failed to get robots list: %v", utils.ParseHarborErrorMsg(err))
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
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
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
