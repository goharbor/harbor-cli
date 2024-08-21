// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// NewGetRegistryCommand creates a new `harbor get registry` command
func ViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "get registry by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registryId, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Errorf("failed to parse registry id: %v", err)
				}
				err = api.GetRegistry(registryId)
				if err != nil {
					log.Errorf("failed to get registry: %v", err)
				}
			} else {
				registryId := prompt.GetRegistryNameFromUser()
				err = api.GetRegistry(registryId)
			}
			if err != nil {
				log.Errorf("failed to get registry: %v", err)
			}
		},
	}

	return cmd
}
