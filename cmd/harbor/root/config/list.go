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
package config

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func ListConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List config items",
		Example: `  harbor config list`,
		Long:    `Get information of all CLI config items`,
		Args:    cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				logrus.Errorf("Failed to get config: %v", err)
				return
			}

			// Get the output format
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				// Use utils.PrintFormat if available
				err = utils.PrintFormat(config, formatFlag)
				if err != nil {
					logrus.Errorf("Failed to print config: %v", err)
				}
			} else {
				// Default to YAML format
				data, err := yaml.Marshal(config)
				if err != nil {
					logrus.Errorf("Failed to marshal config to YAML: %v", err)
					return
				}
				fmt.Println(string(data))
			}
		},
	}

	return cmd
}
