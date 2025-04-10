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

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteInstanceCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete instance by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				instanceName := args[0]
				err = api.DeleteInstance(instanceName)
			} else {
				instanceName := prompt.GetInstanceFromUser()
				err = api.DeleteInstance(instanceName)
			}
			if err != nil {
				log.Errorf("failed to delete instance: %v", err)
			}
		},
	}

	return cmd
}
