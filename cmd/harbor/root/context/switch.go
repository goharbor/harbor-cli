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

package context

import (
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func SwitchContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "switch <none|context>",
		Short:   "Switch to a new context",
		Example: `harbor context switch harbor-cli@https-demo-goharbor-io`,
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				logrus.Errorf("Failed to get config: %v", err)
				return
			}

			if len(args) > 1 {
				logrus.Errorf("Invalid number of arguments")
				return
			} else if len(args) == 1 {
				newActiveCredential := args[0]
				found := false

				for _, cred := range config.Credentials {
					if cred.Name == newActiveCredential {
						found = true
						break
					}
				}
				if found {
					config.CurrentCredentialName = newActiveCredential
					if err := utils.UpdateConfigFile(config); err != nil {
						logrus.Errorf("Failed to update config: %v", err)
					}
				} else {
					logrus.Errorf("Context doesn't exist!")
				}
			} else {
				logrus.Infof("Not implemented yet")
			}
		},
	}
	return cmd
}
