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
package scan_all

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/scan-all/view-schedule"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// This command does not work because the API does not return the response body
// API: https://demo.goharbor.io/devcenter-api-2.0
func ViewScanAllScheduleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view-schedule",
		Short:   "View the scan all schedule",
		Aliases: []string{"vs"},
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Info("Retrieving scan all schedule configuration")
			schedule, err := api.GetScanAllSchedule()
			if err != nil {
				logrus.Errorf("Failed to retrieve scan all schedule: %v", utils.ParseHarborErrorMsg(err))
				return err
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(schedule, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				view.ViewScanSchedule(schedule)
			}
			return nil
		},
	}

	return cmd
}
