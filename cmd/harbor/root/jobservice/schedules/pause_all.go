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
	"github.com/spf13/cobra"
)

// PauseAllCommand pauses all schedules
func PauseAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause-all",
		Short:   "Pause all schedules",
		Long:    "Pause the global scheduler and all schedules.",
		Example: "harbor jobservice schedules pause-all",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Pausing all schedules...")
			err := api.ActionJobQueue("SCHEDULER", "pause")
			if err != nil {
				return formatScheduleError("failed to pause all schedules", err, "ActionStop")
			}
			fmt.Println("✓ All schedules paused successfully.")
			return nil
		},
	}

	return cmd
}
