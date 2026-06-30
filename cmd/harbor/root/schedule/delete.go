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

func DeleteScheduleCommand() *cobra.Command {
	var jobType string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a schedule job (e.g. gc, scan-all)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jobType == "" {
				return fmt.Errorf("--type flag is required")
			}

			err := api.DeleteSchedule(jobType)
			if err != nil {
				return fmt.Errorf("failed to delete %s schedule: %v", jobType, err)
			}

			fmt.Printf("Successfully deleted schedule for %s\n", jobType)
			return nil
		},
	}

	cmd.Flags().StringVar(&jobType, "type", "", "Type of job to delete schedule for (gc, scan-all)")

	return cmd
}
