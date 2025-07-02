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
package scan_all

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func StopScanAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop scanning all artifacts",
		Long: `Stop an ongoing vulnerability scan of all artifacts in Harbor.

This command halts the current scan-all operation that was either manually triggered 
or scheduled. When stopped, scans that are already in progress will complete, but no new artifacts will be scanned. The scan can be restarted later using the 'scan-all run' command.

Examples:
  # Stop the current scan-all operation
  harbor-cli scan-all stop

  # Stop and then check metrics to confirm
  harbor-cli scan-all stop && harbor-cli scan-all metrics`,
		Args: cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Info("Stopping scan all operation")
			err := api.StopScanAll()
			if err != nil {
				logrus.Errorf("Failed to stop scan all operation: %v", utils.ParseHarborErrorMsg(err))
				return err
			}
			logrus.Info("Successfully stopped scan all operation")
			return nil
		},
	}

	return cmd
}
