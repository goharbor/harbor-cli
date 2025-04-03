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
package immutable

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteImmutableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete immutable rule",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var immutableId int64
			if len(args) > 0 {
				immutableId = prompt.GetImmutableTagRule(args[0])
				err = api.DeleteImmutable(args[0], immutableId)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				immutableId = prompt.GetImmutableTagRule(projectName)
				err = api.DeleteImmutable(projectName, immutableId)
			}
			if err != nil {
				log.Errorf("failed to delete immutable rule: %v", err)
			}
		},
	}

	return cmd
}
