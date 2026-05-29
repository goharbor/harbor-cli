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
package instance

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PingInstanceCommand() *cobra.Command {
	var useInstanceID bool

	cmd := &cobra.Command{
		Use:   "ping [NAME|ID]",
		Short: "Ping preheat provider instance by name or id",
		Long: `Ping a preheat provider instance to test its connectivity in Harbor. You can specify the instance
by name or ID directly as an argument. If no argument is provided, you will be prompted to select
an instance from a list of available instances.`,
		Example: `  harbor-cli instance ping my-instance
  harbor-cli instance ping 1 --id`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var instanceName string

			if useInstanceID && len(args) == 0 {
				return fmt.Errorf("instance ID must be provided when using --id")
			}

			if len(args) > 0 {
				log.Debugf("Instance name provided: %s", args[0])
				instanceName = args[0]
			} else {
				log.Debug("No instance name provided, prompting user")
				instanceName, err = prompt.GetInstanceNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get instance name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			log.Debugf("Pinging instance: %s", instanceName)
			response, err := api.PingInstance(instanceName, useInstanceID)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "404" {
					return fmt.Errorf("instance %s does not exist", instanceName)
				}
				return fmt.Errorf("failed to ping instance: %v", utils.ParseHarborErrorMsg(err))
			}

			outputFormat := viper.GetString("output-format")
			if outputFormat != "" {
				if err := utils.PrintFormat(response, outputFormat); err != nil {
					return err
				}
			} else {
				fmt.Printf("Instance '%s' pinged successfully\n", instanceName)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&useInstanceID, "id", false, "Get instance by id")

	return cmd
}
