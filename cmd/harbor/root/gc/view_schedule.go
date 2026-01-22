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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ViewGCScheduleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "schedule",
		Short: "Display the GC schedule",
		Run: func(cmd *cobra.Command, args []string) {
			scheduleWrapper, err := api.GetGCSchedule()
			if err != nil {
				logrus.Fatalf("Failed to get GC schedule: %v", err)
			}

			if scheduleWrapper == nil || scheduleWrapper.Schedule == nil {
				fmt.Println("No GC schedule set.")
				return
			}

			s := scheduleWrapper.Schedule

			fmt.Printf("Schedule Type:     %s\n", s.Type)
			if s.Cron != "" {
				fmt.Printf("Cron Expression:   %s\n", s.Cron)
			}
			fmt.Printf("Next Execution:    %v\n", s.NextScheduledTime)
			fmt.Printf("Creation Time:     %v\n", scheduleWrapper.CreationTime)
			fmt.Printf("Update Time:       %v\n", scheduleWrapper.UpdateTime)
		},
	}
}
