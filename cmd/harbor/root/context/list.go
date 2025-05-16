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
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/context/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List contexts",
		Example: `  harbor context list`,
		Args:    cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				fmt.Errorf("failed to get config: %v", utils.ParseHarborErrorMsg(err))
				return
			}

			// Get the output format
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				// Use utils.PrintFormat if available
				err = utils.PrintFormat(config, formatFlag)
				if err != nil {
					fmt.Errorf("Failed to print config: %v", utils.ParseHarborErrorMsg(err))
					return
				}
			} else {
				var cxlist []api.ContextListView
				for _, cred := range config.Credentials {
					cx := api.ContextListView{Name: cred.Name, Username: cred.Username, Server: cred.ServerAddress}
					cxlist = append(cxlist, cx)
				}
				currentCredential := config.CurrentCredentialName
				list.ListContexts(cxlist, currentCredential)
			}
		},
	}
	return cmd
}
