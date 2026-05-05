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
package policy

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/preheat/policy/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListPolicyCommand() *cobra.Command {
	var opts api.ListFlags
	var isID bool

	cmd := &cobra.Command{
		Use:   "list [NAME|ID]",
		Short: "List preheat policies under a project by name or id",
		Long:  "List project-scoped P2P preheat policies in Harbor",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName string

			if opts.Page < 1 {
				return fmt.Errorf("page number must be greater than or equal to 1")
			}
			if opts.PageSize <= 0 || opts.PageSize > 100 {
				return fmt.Errorf("page size must be greater than 0 and less than or equal to 100")
			}

			if len(args) > 0 {
				log.Debugf("Project name provided: %s", args[0])
				projectName = args[0]
			} else {
				log.Debug("No project name provided, prompting user")
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			log.Debug("Fetching preheat policies...")
			resp, err := api.ListPreheatPolicies(projectName, isID, opts)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "404" {
					return fmt.Errorf("project %s does not exist", projectName)
				}
				return fmt.Errorf("failed to list preheat policies: %v", utils.ParseHarborErrorMsg(err))
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(resp, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				list.ListPolicies(resp.Payload)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&isID, "id", false, "Get preheat policies by project id")
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")

	return cmd
}
