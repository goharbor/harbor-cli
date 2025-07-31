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
package configurations

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update system configurations from local config file",
		Long: `Update Harbor system configurations using the values stored in your local config file.
This will push the configurations from your local config file to the Harbor server.
Make sure to run 'harbor config get' first to populate the local config file with current configurations.`,
		Args:    cobra.NoArgs,
		Example: `harbor config update`,
		RunE: func(cmd *cobra.Command, args []string) error {
			harborConfig, err := utils.GetCurrentHarborConfig()
			if err != nil {
				return fmt.Errorf("failed to get config from file: %v", err)
			}

			if harborConfig.Configurations == (models.Configurations{}) {
				return fmt.Errorf("no configurations found in config file. Run 'harbor config get' first to populate configurations")
			}

			err = api.UpdateConfigurations(harborConfig)
			if err != nil {
				return fmt.Errorf("failed to update Harbor configurations: %v", err)
			}

			log.Infof("Harbor configurations updated successfully from local config file.")
			return nil
		},
	}
	return cmd
}
