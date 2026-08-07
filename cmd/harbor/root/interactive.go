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
	"strings"

	commandinteractive "github.com/goharbor/harbor-cli/pkg/views/command/interactive"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			selected, err := commandinteractive.Run(root)
			if err != nil || selected == nil {
				return err
			}

			return executeInteractiveSelection(cmd.Root(), selected)
		},
	}

	return cmd
}

func executeInteractiveSelection(sourceRoot, selected *cobra.Command) error {
	commandPath := strings.TrimSpace(selected.CommandPath())
	execRoot := RootCmd()
	execRoot.SetArgs(append(changedRootFlags(sourceRoot), commandArgsFromPath(commandPath)...))

	return execRoot.Execute()
}

func commandArgsFromPath(commandPath string) []string {
	parts := strings.Fields(strings.TrimSpace(commandPath))
	if len(parts) <= 1 {
		return nil
	}

	return parts[1:]
}

func changedRootFlags(root *cobra.Command) []string {
	if root == nil {
		return nil
	}

	var args []string
	root.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Changed {
			return
		}

		args = append(args, "--"+flag.Name)
		if flag.Value.Type() != "bool" {
			args = append(args, flag.Value.String())
		}
	})

	return args
}
