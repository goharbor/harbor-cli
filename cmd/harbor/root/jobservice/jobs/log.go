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
package jobs

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	jobserviceutils "github.com/goharbor/harbor-cli/pkg/utils/jobservice"
	"github.com/spf13/cobra"
)

// JobsCommand creates the jobs subcommand
func JobsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Manage job logs (view by job ID)",
		Long:  "View logs for specific jobs.",
	}

	cmd.AddCommand(LogCommand())

	return cmd
}

// LogCommand retrieves and displays job logs
func LogCommand() *cobra.Command {
	var jobID string

	cmd := &cobra.Command{
		Use:     "log",
		Short:   "View a job log (--job-id required)",
		Long:    "Display the log for a specific job by job ID.",
		Example: "harbor jobservice jobs log --job-id abc123def456",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jobID == "" {
				return fmt.Errorf("--job-id must be specified")
			}

			fmt.Printf("Retrieving log for job %s...\n\n", jobID)

			log, err := api.GetJobLog(jobID)
			if err != nil {
				return jobserviceutils.FormatScheduleError("failed to retrieve job log", err, "authenticated")
			}

			if log == "" {
				fmt.Println("No log content available for this job.")
				return nil
			}

			fmt.Println("=== Job Log ===")
			fmt.Println(log)
			fmt.Println("=== End of Log ===")
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&jobID, "job-id", "", "Job ID to fetch log for (required)")
	if err := cmd.MarkFlagRequired("job-id"); err != nil {
		panic(err)
	}

	return cmd
}
