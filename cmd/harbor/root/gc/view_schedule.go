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

package gc

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/gc/schedule"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ViewGCScheduleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "schedule",
		Short: "Display the GC schedule",
		Args:  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			scheduleWrapper, err := api.GetGCSchedule()
			if err != nil {
				return fmt.Errorf("failed to get GC schedule: %v", utils.ParseHarborErrorMsg(err))
			}

			if scheduleWrapper == nil || scheduleWrapper.Schedule == nil {
				log.Info("No GC schedule set.")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(scheduleWrapper, formatFlag)
				if err != nil {
					return err
				}
			} else {
				schedule.ViewGCSchedule(scheduleWrapper)
			}
			return nil
		},
	}
}
