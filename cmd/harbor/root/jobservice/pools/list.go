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
package pools

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	poolviews "github.com/goharbor/harbor-cli/pkg/views/jobservice/pools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PoolsCommand creates the pools subcommand.
func PoolsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Short: "Manage worker pools",
		Long:  "List and manage worker pools for the Harbor job service.",
	}

	cmd.AddCommand(ListCommand())
	return cmd
}

// ListCommand lists all worker pools.
func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all worker pools",
		Long:    "Display all worker pools with their details.",
		Example: "harbor jobservice pools list",
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := api.GetWorkerPools()
			if err != nil {
				return fmt.Errorf("failed to retrieve worker pools: %w", err)
			}

			if response == nil || response.Payload == nil || len(response.Payload) == 0 {
				fmt.Println("No worker pools found.")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				return utils.PrintFormat(response.Payload, formatFlag)
			}

			poolviews.ListPools(response.Payload)
			return nil
		},
	}

	return cmd
}
