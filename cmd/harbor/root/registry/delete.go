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
package registry

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteRegistryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete registry",
		Example: "harbor registry delete [registryname]",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registryName, _ := api.GetRegistryIdByName(args[0])
				err = api.DeleteRegistry(registryName)
			} else {
				registryId := prompt.GetRegistryNameFromUser()
				err = api.DeleteRegistry(registryId)
			}
			if err != nil {
				log.Errorf("failed to delete registry: %v", err)
			}
		},
	}

	return cmd
}
