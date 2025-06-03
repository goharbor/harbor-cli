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
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
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
				fmt.Println("failed to get config: ", utils.ParseHarborErrorMsg(err))
				return
			}

			if len(args) == 1 {
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
						fmt.Println("failed to update config: ", utils.ParseHarborErrorMsg(err))
					}
				} else {
					fmt.Println("context doesn't exist")
				}
			} else {
				res, err := prompt.GetActiveContextFromUser()
				if err != nil {
					fmt.Println("failed to get active context: ", utils.ParseHarborErrorMsg(err))
					return
				}
				if res != "" {
					msg := fmt.Sprintf("context switched from '%s' to '%s'", config.CurrentCredentialName, res)
					config.CurrentCredentialName = res
					if err := utils.UpdateConfigFile(config); err != nil {
						fmt.Println("failed to update config: ", utils.ParseHarborErrorMsg(err))
					} else {
						fmt.Println(msg)
					}
				}
			}
		},
	}
	return cmd
}
