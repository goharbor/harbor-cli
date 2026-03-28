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
	workersviews "github.com/goharbor/harbor-cli/pkg/views/jobservice/workers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outputFormat string

// ListCommand lists all workers
func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [POOL_ID]",
		Short: "List job service workers",
		Long: `List job service workers for a specific worker pool.

If no pool ID is specified, it will use 'all' to get workers from all pools.

Examples:
  harbor jobservice workers list                # List workers from all pools
  harbor jobservice workers list default        # List workers from 'default' pool
  harbor jobservice workers list my-pool        # List workers from 'my-pool'`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listWorkers(cmd, args)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table",
		"output format (table, json, yaml)")

	return cmd
}

func listWorkers(cmd *cobra.Command, args []string) error {
	poolID := "all" // default to all pools
	if len(args) > 0 {
		poolID = args[0]
	}

	resp, err := api.GetWorkers(poolID)
	if err != nil {
		return fmt.Errorf("failed to get workers: %w", err)
	}

	if resp == nil || resp.Payload == nil || len(resp.Payload) == 0 {
		fmt.Println("No workers found.")
		return nil
	}

	formatFlag := viper.GetString("output-format")
	if formatFlag != "" {
		return utils.PrintFormat(resp.Payload, formatFlag)
	}

	workersviews.ListWorkers(resp.Payload)
	return nil
}
