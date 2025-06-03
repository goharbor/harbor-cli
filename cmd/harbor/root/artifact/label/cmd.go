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
package label

import (
	"github.com/spf13/cobra"
)

// LabelsArtifactCommmand compound command to label artifacts
func LabelsArtifactCommmand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "label command for artifacts",
		Long:  `label command for artifact`,
		Example: `harbor artifact label add <project>/<repository>/<reference> <label name>
harbor artifact label del <project>/<repository>/<reference> <label name>
		`,
	}
	cmd.AddCommand(AddLabelArtifactCommmand())
	cmd.AddCommand(DelLabelArtifactCommmand())
	cmd.AddCommand(ListLabelArtifactCommmand())
	return cmd
}
