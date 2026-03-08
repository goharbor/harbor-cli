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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RunGCCommand() *cobra.Command {
	var dryRun, deleteUntagged bool

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run Garbage Collection manually",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			scheduleObj := models.ScheduleObj{
				Type: "Manual",
			}

			params := map[string]interface{}{
				"dry_run":         dryRun,
				"delete_untagged": deleteUntagged,
			}

			scheduleBody := &models.Schedule{
				Schedule:   &scheduleObj,
				Parameters: params,
			}

			err := api.CreateGCSchedule(scheduleBody)
			if err != nil {
				return fmt.Errorf("failed to start GC: %v", utils.ParseHarborErrorMsg(err))
			}
			log.Info("GC started successfully")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Simulate GC without deleting artifacts")
	cmd.Flags().BoolVarP(&deleteUntagged, "delete-untagged", "", true, "Delete untagged artifacts")

	return cmd
}
