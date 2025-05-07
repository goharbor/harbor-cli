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
package context

import "github.com/spf13/cobra"

func Context() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "context",
		Short:   "Manage locally available contexts",
		Example: "harbor context list",
		Long: `The context command allows you to manage configs of the Harbor CLI.
				You can add, get, or delete specific config item, as well as list all config items of the Harbor Cli`,
	}

	cmd.AddCommand(
		ListContextCommand(),
		GetContextItemCommand(),
		UpdateContextItemCommand(),
		DeleteContextItemCommand(),
	)

	return cmd
}
