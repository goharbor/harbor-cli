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

// ListCommand lists all workers
func ListCommand() *cobra.Command {
	var poolID string
	var allPools bool
	var poolAll bool

	cmd := &cobra.Command{
		Use:   "list [POOL_ID]",
		Short: "List workers (all pools by default; use --pool for one pool)",
		Long: `List job service workers.

Supported listing modes:
	- All workers (default): no POOL_ID or --pool all
	- Specific pool workers: provide [POOL_ID] or --pool <pool-id>
	- Compatibility mode: --pool-all (same as --pool all)

Examples:
  harbor jobservice workers list
  harbor jobservice workers list --pool all
  harbor jobservice workers list --pool default
	harbor jobservice workers list default
	harbor jobservice worker list 72327cf790564e45b7c89a2d`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listWorkers(cmd, args, poolID, poolAll)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&poolID, "pool", "", "Worker pool ID (use 'all' for all pools)")
	flags.BoolVar(&allPools, "all", false, "List workers from all pools")
	flags.BoolVar(&poolAll, "pool-all", false, "List workers from all pools (compatibility alias for --pool all)")

	return cmd
}

func listWorkers(cmd *cobra.Command, args []string, poolFlag string, poolAll bool) error {
	resolvedPoolID := "all"
	allPools, _ := cmd.Flags().GetBool("all")

	if allPools || poolAll {
		resolvedPoolID = "all"
	}

	if poolFlag != "" {
		resolvedPoolID = poolFlag
	}

	if len(args) > 0 {
		if poolFlag != "" || allPools || poolAll {
			return fmt.Errorf("pool ID provided both as argument and flag; use only one form")
		}
		resolvedPoolID = args[0]
	}

	resp, err := api.GetWorkers(resolvedPoolID)
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
