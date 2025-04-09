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
package retention

import (
	"github.com/spf13/cobra"
)

func Retention() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retention",
		Short: "Manage tag retention rules in the project",
		Long: `Manage tag retention rules in the project in Harbor.
		
The 'retention' command allows users to create, list, and delete tag retention rules 
within a project. Tag retention rules help in managing and controlling the lifecycle 
of tags by defining rules for automatic cleanup and retention.

A user can only create up to 15 tag retention rules per project.`,
		Example: `  harbor tag retention create    # Create a new tag retention rule
  harbor tag retention list      # List all tag retention rules in the project
  harbor tag retention delete    # Delete a specific tag retention rules`,
	}

	cmd.AddCommand(
		CreateRetentionCommand(),
		ListRetentionRulesCommand(),
		DeleteRetentionRuleCommand(),
	)
	return cmd
}
