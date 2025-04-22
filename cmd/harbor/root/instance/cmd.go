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
package instance

import "github.com/spf13/cobra"

func Instance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Manage preheat provider instances in Harbor",
		Long: `Manage preheat provider instances used by Harbor for pre-distributing container images.
These instances represent external services such as Dragonfly or Kraken that help preheat images across nodes.`,
	}
	cmd.AddCommand(
		CreateInstanceCommand(),
		DeleteInstanceCommand(),
		ListInstanceCommand(),
	)
	return cmd
}
