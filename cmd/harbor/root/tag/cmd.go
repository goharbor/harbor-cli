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
package tag

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/tag/immutable"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/tag/retention"
	"github.com/spf13/cobra"
)

func TagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tag rules in Harbor registry",
		Long:  "Manage rules for tags in the Harbor registry, including creating, listing, and deleting them.",
	}
	cmd.AddCommand(retention.Retention())
	cmd.AddCommand(immutable.Immutable())
	return cmd
}
