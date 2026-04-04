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
package schedules

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/jobservice/schedules"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StatusCommand shows the global scheduler status
func StatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Show scheduler status",
		Long:    "Display whether the global scheduler is paused or running.",
		Example: "harbor jobservice schedules status",
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := api.GetSchedulePaused()
			if err != nil {
				return formatScheduleError("failed to retrieve scheduler status", err, "authenticated")
			}

			if response == nil || response.Payload == nil {
				fmt.Println("Unable to determine scheduler status.")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				return utils.PrintFormat(response.Payload, formatFlag)
			}

			schedules.PrintScheduleStatus(response.Payload)
			return nil
		},
	}

	return cmd
}
