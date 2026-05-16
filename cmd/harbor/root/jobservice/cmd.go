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

package jobservice

import (
	"github.com/spf13/cobra"
)

func JobServiceCmd() *cobra.Command {
	// jobserviceCmd represents the jobservice command.
	var jobserviceCmd = &cobra.Command{
		Use:     "jobservice",
		Aliases: []string{"js"},
		Short:   "Manage Harbor job service",
		Long: `Manage Harbor job service, including queues, jobs, worker pools and workers.

This command provides terminal-based access to the Jobservice dashboard, allowing you to:
- Monitor and manage job queues (pause, resume, clear)
- View running jobs and stop them if necessary
- Inspect worker pools and active workers
- Access job logs with real-time tailing support`,
		Example: `  # List all job queues
  harbor jobservice queue list

  # Pause a specific job queue
  harbor jobservice queue pause IMAGE_SCAN

  # Stop a running job
  harbor jobservice job stop <job-id>

  # Follow job logs in real-time
  harbor jobservice job log <job-id> --follow`,
	}

	jobserviceCmd.AddCommand(
		QueueCommand(),
		JobCommand(),
		PoolCommand(),
		WorkerCommand(),
	)

	return jobserviceCmd
}
