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
	"github.com/goharbor/harbor-cli/pkg/utils"
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
				return formatWorkerActionError("failed to free worker", err)
			}

			fmt.Printf("Worker job %q stopped successfully.\n", jobID)
			return nil
		},
	}

	cmd.Flags().StringVar(&jobID, "job-id", "", "Running job ID to stop")

	return cmd
}

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
				return formatWorkerActionError("failed to free all workers", err)
			}

			fmt.Println("All busy workers were freed successfully.")
			return nil
		},
	}

	return cmd
}

func formatWorkerActionError(operation string, err error) error {
	errorCode := utils.ParseHarborErrorCode(err)

	switch errorCode {
	case "401":
		return fmt.Errorf("%s: authentication required. Please run 'harbor login' and try again", operation)
	case "403":
		return fmt.Errorf("%s: permission denied. This operation requires ActionStop on jobservice-monitor", operation)
	case "404":
		return fmt.Errorf("%s: job not found or already completed", operation)
	case "500":
		return fmt.Errorf("%s: Harbor internal error. Retry and check Harbor server logs", operation)
	default:
		msg := utils.ParseHarborErrorMsg(err)
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("%s: %s", operation, msg)
	}
}
