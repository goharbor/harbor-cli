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
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/health"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func HealthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Get the health status of Harbor components",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := api.Ping()
			if err != nil {
				return err
			}
			status, err := api.GetHealth()
			if err != nil {
				return err
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(status, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				health.PrintHealthStatus(status)
			}
			return nil
		},
		Example: `  # Get the health status of Harbor components`,
	}

	return cmd
}
