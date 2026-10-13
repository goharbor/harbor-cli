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
	"github.com/spf13/cobra"
)

// PoolsCommand creates the pools subcommand.
func PoolsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pools",
		Aliases: []string{"pool"},
		Short:   "Manage worker pools (list available pools)",
		Long: `List and manage worker pools for the Harbor job service.

Use 'list' to view all worker pools.

Examples:
  harbor jobservice pools list
  harbor jobservice pool list`,
	}

	cmd.AddCommand(ListCommand())
	return cmd
}
