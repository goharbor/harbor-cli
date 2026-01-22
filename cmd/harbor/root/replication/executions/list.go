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
package executions

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/replication/execution/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCommand() *cobra.Command {
	var opts api.ListFlags
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List replication executions",
		Long:  `List all replication executions for a given replication policy. If no policy ID is provided, it will prompt the user to select one interactively.`,
		Example: `  harbor replication executions list 12345
  harbor replication executions list`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("Starting replication executions list command")

			var rpolicyID int64
			if len(args) > 0 {
				var err error
				// convert string to int64
				rpolicyID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid replication policy ID: %s, %v", args[0], err)
				}
			} else {
				rpolicyID = prompt.GetReplicationPolicyFromUser()
			}

			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}

			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

			log.Debug("Fetching executions...")
			executions, err := api.ListReplicationExecutions(rpolicyID, opts)
			if err != nil {
				return fmt.Errorf("failed to get projects list: %v", utils.ParseHarborErrorMsg(err))
			}

			log.WithField("count", len(executions.Payload)).Debug("Number of executions fetched")
			if len(executions.Payload) == 0 {
				log.Info("No executions found")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				log.WithField("output_format", formatFlag).Debug("Output format selected")
				err = utils.PrintFormat(executions.Payload, formatFlag)
				if err != nil {
					return err
				}
			} else {
				log.Debug("Listing projects using default view")
				list.ListExecutions(executions.Payload)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 0, "Size of per page (0 to fetch all)")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")

	return cmd
}
