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

// FreeCommand frees a worker by stopping the running job on it.
func FreeCommand() *cobra.Command {
	var jobID string

	cmd := &cobra.Command{
		Use:     "free",
		Short:   "Free one worker (--job-id required)",
		Long:    "Stop a running job by job ID to free its worker.",
		Example: "harbor jobservice workers free --job-id abc123",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jobID == "" {
				return fmt.Errorf("--job-id is required")
			}

			err := api.StopRunningJob(jobID)
			if err != nil {
				return jobserviceutils.FormatScheduleError("failed to free worker", err, "ActionStop")
			}

			fmt.Printf("Worker job %q stopped successfully.\n", jobID)
			return nil
		},
	}

	cmd.Flags().StringVar(&jobID, "job-id", "", "Running job ID to stop")

	return cmd
}
