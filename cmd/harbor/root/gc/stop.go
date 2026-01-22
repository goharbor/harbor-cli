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

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/gc/stop"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func StopGCCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a running GC job",
		Long: `Stop a running Garbage Collection job in Harbor.

This command displays a list of currently running or pending GC jobs and 
allows you to select one to stop. Only jobs with status "running" or "pending" 
can be stopped.

Examples:
  # Stop a running GC job interactively
  harbor-cli gc stop

Notes:
  - Only jobs that are currently running or pending can be stopped
  - Jobs that have already completed cannot be stopped
  - Use 'harbor-cli gc list' to view all GC jobs and their statuses`,
		Args: cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			gcID, err := stop.SelectGCJob()
			if err != nil {
				if err.Error() == "no running GC jobs found to stop" {
					logrus.Info(err.Error())
					return nil
				}
				logrus.Errorf("Failed to select GC job: %v", err)
				return err
			}

			logrus.Infof("Stopping GC job %d", gcID)
			err = api.StopGC(gcID)
			if err != nil {
				errMsg := utils.ParseHarborErrorMsg(err)
				if contains(errMsg, "404") {
					return fmt.Errorf("GC job %d not found or already completed", gcID)
				}
				if contains(errMsg, "400") {
					return fmt.Errorf("GC job %d cannot be stopped (may have already completed)", gcID)
				}
				return fmt.Errorf("failed to stop GC job %d: %v", gcID, errMsg)
			}

			logrus.Infof("Successfully stopped GC job %d", gcID)
			return nil
		},
	}

	return cmd
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
