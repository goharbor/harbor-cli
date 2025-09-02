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
		Example: `harbor config apply -f <config_file>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var configurations *models.Configurations
			var err error
			if cfgFile != "" {
				data, err := os.ReadFile(cfgFile)
				if err != nil {
					return fmt.Errorf("failed to read config file: %v", err)
				}
				fileType := filepath.Ext(cfgFile)
				switch fileType {
				case ".yaml", ".yml":
					if err := yaml.Unmarshal(data, &configurations); err != nil {
						return fmt.Errorf("failed to parse YAML: %v", err)
					}
				case ".json":
					if err := json.Unmarshal(data, &configurations); err != nil {
						return fmt.Errorf("failed to parse JSON: %v", err)
					}
				default:
					return fmt.Errorf("unsupported file type: %s, expected '.yaml/.yml' or '.json'", fileType)
				}
			} else {
				return fmt.Errorf("no config file specified")
			}

			err = api.UpdateConfigurations(configurations)
			if err != nil {
				return fmt.Errorf("failed to update Harbor configurations: %v", err)
			}

			fmt.Printf("harbor configurations updated successfully from %s.", cfgFile)
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&cfgFile, "configurations-file", "f", "", "Harbor configurations file to apply.")

	return cmd
}
