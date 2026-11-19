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
	"github.com/spf13/cobra"
)

// WorkersCommand creates the workers subcommand.
func WorkersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workers",
		Aliases: []string{"worker"},
		Short:   "Manage workers (list all/by pool, free, free-all)",
		Long: `Manage job service workers using the job service API.

Use 'list' to view workers from all pools or a specific pool.
Use 'free' and 'free-all' to stop running jobs and release busy workers.

Examples:
  harbor jobservice workers list
  harbor jobservice workers list --pool all
  harbor jobservice workers list --pool default
	harbor jobservice workers list --page 2 --page-size 20
	harbor jobservice workers list default
	harbor jobservice worker list 72327cf790564e45b7c89a2d
  harbor jobservice workers free --job-id <job-id>
  harbor jobservice workers free-all`,
	}

	cmd.AddCommand(
		ListCommand(),
		FreeCommand(),
		FreeAllCommand(),
	)

	return cmd
}
