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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func ApplyConfigCmd() *cobra.Command {
	var cfgFile string
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Update system configurations from local config file",
		Long: `Update Harbor system configurations using the values stored in your local config file.
This will push the configurations from your local config file to the Harbor server.
Make sure to run 'harbor config get' first to populate the local config file with current configurations. Alternatively, you can specify a custom configuration file using the --configurations-file flag. This does not have to be a complete configuration file, only the fields you want to update need to be present under the 'configurations' key. Credentials for the Harbor server can be configured in the local config file or through environment variables or global config flags.`,
		Args:    cobra.NoArgs,
		Example: `harbor config apply`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var harborConfig *utils.HarborConfig
			var err error
			if cfgFile != "" {
				data, err := os.ReadFile(cfgFile)
				if err != nil {
					return fmt.Errorf("failed to read config file: %v", err)
				}
				fileType := filepath.Ext(cfgFile)
				switch fileType {
				case ".yaml", ".yml":
					if err := yaml.Unmarshal(data, &harborConfig); err != nil {
						return fmt.Errorf("failed to parse YAML: %v", err)
					}
				case ".json":
					if err := json.Unmarshal(data, &harborConfig); err != nil {
						return fmt.Errorf("failed to parse JSON: %v", err)
					}
				default:
					return fmt.Errorf("unsupported file type: %s, expected '.yaml/.yml' or '.json'", fileType)
				}
			} else {
				harborConfig, err = utils.GetCurrentHarborConfig()
				if err != nil {
					return fmt.Errorf("failed to get config from file: %v", err)
				}
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
	flags := cmd.Flags()
	flags.StringVarP(&cfgFile, "configurations-file", "f", "", "Harbor configurations file to apply (default is $HOME/.harbor/config.yaml). This file should contain the 'configurations' key with the fields you want to update. If not specified, it will use the default Harbor config file.")

	return cmd
}
