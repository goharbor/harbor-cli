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
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/registry/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ViewRegistryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "get registry information",
		Example: "harbor registry view [registryName]",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var registryId int64
			var registry *registry.GetRegistryOK

			if len(args) > 0 {
				registryId, err = api.GetRegistryIdByName(args[0])
				if err != nil {
					return fmt.Errorf("failed to get registry id: %v", utils.ParseHarborError(err))
				}
			} else {
				registryId = prompt.GetRegistryNameFromUser()
			}

			registry, err = api.ViewRegistry(registryId)
			if err != nil {
				return fmt.Errorf("failed to get registry info: %v", utils.ParseHarborError(err))
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(registry, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				view.ViewRegistry(registry.Payload)
			}
			return nil
		},
	}

	return cmd
}
