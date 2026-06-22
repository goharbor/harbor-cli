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
	"github.com/goharbor/harbor-cli/pkg/views/gc/update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func TriggerGCOperation() *cobra.Command {
	var deleteUntagged bool
	var dryRun bool
	var interactive bool

	cmd := &cobra.Command{
		Use:     "trigger",
		Short:   "Trigger Garbage Collection immediately",
		Long:    `Start a manual Garbage Collection job immediately in Harbor registry.`,
		Example: `  harbor gc trigger --delete-untagged --dry-run=false`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactive {
				logrus.Debug("Opening interactive form for triggering GC")
				update.EditGCSchedule(nil, &deleteUntagged, &dryRun, false)
			}

			logrus.Debugf("Triggering manual GC (delete_untagged: %t, dry_run: %t)", deleteUntagged, dryRun)
			err := api.TriggerGC(deleteUntagged, dryRun)
			if err != nil {
				return fmt.Errorf("failed to trigger GC: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Println("Garbage Collection triggered successfully")
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&deleteUntagged, "delete-untagged", false, "Delete untagged artifacts")
	flags.BoolVar(&dryRun, "dry-run", false, "Simulate the GC process without deleting actual blobs")
	flags.BoolVarP(&interactive, "interactive", "i", false, "Trigger Garbage Collection interactively")

	return cmd
}
