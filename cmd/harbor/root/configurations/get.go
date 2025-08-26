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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func GetConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Harbor configurations",
		Long: `Get Harbor system configurations.
		
This command retrieves the current configurations from Harbor and stores them in your local config file.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := api.GetConfigurations()
			if err != nil {
				return err
			}
			if err := utils.AddConfigurationsToConfigFile(response.Payload); err != nil {
				return fmt.Errorf("failed to update config file: %v", err)
			}
			data, err := utils.GetCurrentHarborData()
			fmt.Printf("✓ Configurations have been added to the config file.\n")
			fmt.Printf("Config file location: %s\n", data.ConfigPath)
			return nil
		},
	}
	return cmd
}
