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

import "github.com/spf13/cobra"

func GC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gc",
		Short: "Manage Garbage Collection in Harbor",
		Long: `Use this command to manage registry-wide Garbage Collection (GC) in your Harbor instance.

Garbage Collection cleans up deleted or orphaned blobs/tags in the registry to free up storage space.
This command supports listing execution history, viewing logs, showing schedule configuration, stopping running jobs, and triggering manual runs.`,
		Example: `  # View Garbage Collection execution history
  harbor gc history

  # Get the current Garbage Collection schedule
  harbor gc schedule

  # Trigger Garbage Collection run immediately
  harbor gc trigger --delete-untagged --dry-run=false

  # View execution logs for a GC run
  harbor gc log 12

  # Stop a running Garbage Collection run
  harbor gc stop 12

  # Update the automatic Garbage Collection schedule
  harbor gc update-schedule daily --delete-untagged`,
	}

	cmd.AddCommand(
		HistoryGCOperation(),
		ScheduleGCOperation(),
		TriggerGCOperation(),
		LogGCOperation(),
		StopGCOperation(),
		UpdateGCScheduleCommand(),
	)

	return cmd
}
