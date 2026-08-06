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

package workers

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	jobserviceutils "github.com/goharbor/harbor-cli/pkg/utils/jobservice"
	"github.com/spf13/cobra"
)

// FreeAllCommand frees all busy workers by stopping all running jobs.
func FreeAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "free-all",
		Short:   "Free all busy workers (job-id=all)",
		Long:    "Stop all running jobs to free all busy workers.",
		Example: "harbor jobservice workers free-all",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := api.StopRunningJob("all")
			if err != nil {
				return jobserviceutils.FormatScheduleError("failed to free all workers", err, "ActionStop")
			}

			fmt.Println("All busy workers were freed successfully.")
			return nil
		},
	}

	return cmd
}
