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
package preheat

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project/preheat/execution"
	"github.com/spf13/cobra"
)

func ExecutionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "execution",
		Aliases: []string{"exec"},
		Short:   "Manage preheat executions",
		Long:    "Manage P2P preheat executions under a project",
		Example: `  harbor-cli project preheat execution list [NAME|ID] [POLICY_NAME]`,
	}

	cmd.AddCommand(
		execution.ListExecutionCommand(),
	)

	return cmd
}
