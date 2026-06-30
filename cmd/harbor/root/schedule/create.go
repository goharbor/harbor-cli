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

package schedule

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

func CreateScheduleCommand() *cobra.Command {
	var jobType string
	var cronString string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a schedule job (e.g. gc, scan-all)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jobType == "" {
				return fmt.Errorf("--type flag is required")
			}
			if cronString == "" {
				return fmt.Errorf("--cron flag is required")
			}

			err := api.CreateSchedule(jobType, cronString)
			if err != nil {
				return fmt.Errorf("failed to create %s schedule: %v", jobType, err)
			}

			fmt.Printf("Successfully created schedule for %s with cron: %s\n", jobType, cronString)
			return nil
		},
	}

	cmd.Flags().StringVar(&jobType, "type", "", "Type of job (gc, scan-all)")
	cmd.Flags().StringVar(&cronString, "cron", "", "Cron string for the schedule (e.g. '0 0 * * *')")

	return cmd
}
