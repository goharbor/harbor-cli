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
package root

import (
	commandinteractive "github.com/goharbor/harbor-cli/pkg/views/command/interactive"
	"github.com/spf13/cobra"
)

func InteractiveCommand(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interactive",
		Short: "Browse the Harbor command tree interactively",
		Long: `Browse Harbor commands and subcommands through an interactive menu.

Version 1 focuses on discovery: you can navigate the full command tree, inspect
leaf commands, and review the exact command path before running anything manually.`,
		Example: `  harbor interactive`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := commandinteractive.Run(root)
			return err
		},
	}

	return cmd
}
