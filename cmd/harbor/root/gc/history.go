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
package gc

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	historyView "github.com/goharbor/harbor-cli/pkg/views/gc/history"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GCHistoryCommand() *cobra.Command {
	var (
		opts api.ListFlags
		// For querying, opts.Q
		fuzzy  []string
		match  []string
		ranges []string
	)

	cmd := &cobra.Command{
		Use:   "history",
		Short: "List cleanup jobs in a Harbor instance",
		Long: `List all cleanup jobs in a Harbor instance. 
Displays key metadata including job kind, schedule, creation time and more.

Examples:
  harbor-cli cleanup list                # Interactive prompt for project and repository

Supports pagination, search queries, and sorting using flags.`,

		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

			if len(fuzzy) != 0 || len(match) != 0 || len(ranges) != 0 { // Only Building Query if a param exists
				q, qErr := utils.BuildQueryParam(fuzzy, match, ranges,
					[]string{"job_name", "id", "job_kind", "creation_time", "update_time", "job_status", "deleted"},
				)
				if qErr != nil {
					return qErr
				}

				opts.Q = q
			}

			resp, err := api.ListGCs(&opts)
			if err != nil {
				return err
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(resp.Payload, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				historyView.GCHistory(resp.Payload)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "p", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "n", 10, "Size of per page")
	flags.StringVarP(&opts.Sort, "sort", "s", "", "Sort the resource list in ascending or descending order")
	flags.StringSliceVar(&fuzzy, "fuzzy", nil, "Fuzzy match filter (key=value)")
	flags.StringSliceVar(&match, "match", nil, "exact match filter (key=value)")
	flags.StringSliceVar(&ranges, "range", nil, "range filter (key=min~max)")

	return cmd
}
