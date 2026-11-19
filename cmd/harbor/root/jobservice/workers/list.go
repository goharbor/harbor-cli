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

package workers

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	jobserviceutils "github.com/goharbor/harbor-cli/pkg/utils/jobservice"
	workersviews "github.com/goharbor/harbor-cli/pkg/views/jobservice/workers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListCommand lists all workers
func ListCommand() *cobra.Command {
	var poolID string
	var page int64 = 1
	var pageSize int64 = 20

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workers (supports --page and --page-size)",
		Long: `List job service workers.

Pagination:
	- --page selects the 1-based page number
	- --page-size controls how many workers are shown per page

Examples:
  harbor jobservice workers list
  harbor jobservice workers list --pool default
	harbor jobservice workers list --page 2 --page-size 20`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listWorkers(poolID, page, pageSize)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&poolID, "pool", "all", "Worker pool ID to list workers from (default: all)")
	flags.Int64Var(&page, "page", 1, "Page number")
	flags.Int64Var(&pageSize, "page-size", 20, "Number of workers per page")

	return cmd
}

func listWorkers(poolID string, page, pageSize int64) error {
	if page < 1 {
		return fmt.Errorf("page must be >= 1")
	}
	if pageSize < 1 {
		return fmt.Errorf("page-size must be >= 1")
	}

	resp, err := api.GetWorkers(poolID)
	if err != nil {
		return jobserviceutils.FormatScheduleError("failed to get workers", err, "ActionList")
	}

	if resp == nil || resp.Payload == nil || len(resp.Payload) == 0 {
		fmt.Println("No workers found.")
		return nil
	}

	totalCount := int64(len(resp.Payload))
	start := (page - 1) * pageSize
	if start >= totalCount {
		fmt.Println("No workers found.")
		return nil
	}

	end := start + pageSize
	if end > totalCount {
		end = totalCount
	}

	pageWorkers := resp.Payload[int(start):int(end)]

	formatFlag := viper.GetString("output-format")
	if formatFlag != "" {
		return utils.PrintFormat(pageWorkers, formatFlag)
	}

	workersviews.ListWorkers(pageWorkers, page, pageSize, totalCount)
	return nil
}
