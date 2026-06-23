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

package queues

import (
	"fmt"

	"github.com/spf13/cobra"
)

// PauseCommand pauses a job queue
func PauseCommand() *cobra.Command {
	var jobTypes []string
	var interactive bool

	cmd := &cobra.Command{
		Use:     "pause",
		Short:   "Pause queue(s) (--type or --interactive)",
		Long:    "Pause a job queue or all queues.",
		Example: "harbor jobservice queues pause --type REPLICATION\nharbor jobservice queues pause --type REPLICATION --type RETENTION\nharbor jobservice queues pause --type all",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(jobTypes) == 0 && !interactive {
				interactive = true
			}

			if interactive {
				selectedTypes, err := selectQueueTypes("pause")
				if err != nil {
					return err
				}
				jobTypes = selectedTypes
			}

			if len(jobTypes) == 0 {
				return fmt.Errorf("at least one job type must be specified with --type or interactive mode")
			}

			return executeQueueAction("pause", jobTypes)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVar(&jobTypes, "type", nil, "Job type(s) to pause (repeat flag or comma-separate values; use 'all' for all queues)")
	flags.BoolVarP(&interactive, "interactive", "i", false, "Interactive mode to choose queue type(s) instead of passing --type")

	return cmd
}
