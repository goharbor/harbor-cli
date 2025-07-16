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
package replication

import (
	executions "github.com/goharbor/harbor-cli/cmd/harbor/root/replication/executions"
	"github.com/spf13/cobra"
)

func ReplicationExecutionsCommand() *cobra.Command {
	// replicationCmd represents the replication command.
	var replicationCmd = &cobra.Command{
		Use:     "executions",
		Aliases: []string{"exec"},
		Short:   "Manage replication executions",
		Long:    `Manage replication executions in Harbor context`,
	}
	replicationCmd.AddCommand(
		executions.ListCommand(),
		executions.ViewCommand(),
	)

	return replicationCmd
}
