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
	"github.com/goharbor/harbor-cli/pkg/views/gc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var validGCSortFields = []string{
	"creation_time",
	"update_time",
	"id",
	"job_status",
}

var validGCQueryKeys = []string{
	"id",
	"job_status",
	"job_name",
}

func ListGCCommand() *cobra.Command {
	var (
		opts   api.ListFlags
		sort   []string
		fuzzy  []string
		match  []string
		ranges []string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List GC history",
		Long: `List GC (Garbage Collection) history in Harbor.

This command displays a list of GC executions with their status, creation time,
and other details. You can control the output using pagination flags and format options.

Examples:
  # List GC history with default pagination (page 1, 10 items per page)
  harbor gc list

  # List GC history with custom pagination
  harbor gc list --page 2 --page-size 20

  # List GC history with sorting by creation time (newest first)
  harbor gc list --sort -creation_time

  # List GC history with multiple sort fields
  harbor gc list --sort creation_time --sort -job_status

  # Filter GC history by status (exact match)
  harbor gc list --match job_status=Success

  # Filter GC history by fuzzy match
  harbor gc list --fuzzy job_name=gc`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

			if len(sort) > 0 {
				sortParam, err := utils.BuildSortParam(sort, validGCSortFields)
				if err != nil {
					return err
				}
				opts.Sort = sortParam
			}

			if len(fuzzy) != 0 || len(match) != 0 || len(ranges) != 0 {
				q, qErr := utils.BuildQueryParam(fuzzy, match, ranges, validGCQueryKeys)
				if qErr != nil {
					return qErr
				}
				opts.Q = q
			}

			history, err := api.GetGCHistory(opts)
			if err != nil {
				return fmt.Errorf("failed to get GC history: %v", utils.ParseHarborErrorMsg(err))
			}

			if len(history) == 0 {
				log.Info("No GC history found")
				return nil
			}

			gc.ListGC(history)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "p", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "s", 10, "Size of per page")
	flags.StringSliceVar(&sort, "sort", nil, "Sort the resource list (e.g. --sort creation_time --sort -update_time)")
	flags.StringSliceVar(&fuzzy, "fuzzy", nil, "Fuzzy match filter (key=value)")
	flags.StringSliceVar(&match, "match", nil, "Exact match filter (key=value)")
	flags.StringSliceVar(&ranges, "range", nil, "Range filter (key=min~max)")

	return cmd
}
