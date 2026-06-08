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
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project/preheat/task"
	"github.com/spf13/cobra"
)

func TaskCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "task",
		Short:   "Manage preheat tasks",
		Long:    "Manage related tasks for the given preheat execution",
		Example: `  harbor-cli project preheat task list [PROJECT_NAME|ID] [POLICY_NAME] [EXECUTION_ID]`,
	}

	cmd.AddCommand(
		task.ListTaskCommand(),
		task.LogTaskCommand(),
	)

	return cmd
}
