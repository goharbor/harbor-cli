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
package root

import (
	"fmt"

	"github.com/goharbor/harbor-cli/cmd/harbor/internal/version"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/info/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Lists the info of the Harbor system
func InfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Display detailed Harbor system, statistics, and CLI environment information",
		Long: `The 'info' command retrieves and displays general information about the Harbor instance, 
including system metadata, storage statistics, and CLI environment details such as user identity, 
registry address, and CLI version.

The output can be formatted as table (default), JSON, or YAML using the '--output-format' flag.`,
		Example: `  harbor info
  harbor info --output-format json
  harbor info -o yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var cliinfo *api.CLIInfo
			var err error
			generalInfo, err := api.GetSystemInfo()
			if err != nil {
				return err
			}

			stats, err := api.GetStats()
			if err != nil {
				return err
			}

			sysVolume, err := api.GetSystemVolumes()
			if err != nil {
				return err
			}
			cliVersion := version.Version
			OSinfo := version.System

			cliinfo, err = api.GetCLIInfo()
			if err != nil {
				return fmt.Errorf("Failed to get CLI info: %w", err)
			}
			systemInfo := list.CreateSystemInfo(
				generalInfo.Payload,
				stats.Payload,
				sysVolume.Payload,
				cliinfo,
				cliVersion,
				OSinfo,
			)

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(systemInfo, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.ListInfo(&systemInfo)
			}

			return nil
		},
	}
	return cmd
}
