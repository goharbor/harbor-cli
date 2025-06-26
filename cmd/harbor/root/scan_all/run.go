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
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RunScanAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Scan all artifacts now",
		Long: `Initiate an immediate vulnerability scan of all artifacts in Harbor.

This command triggers a manual scan of all artifacts across all projects in your Harbor instance.
The scan will check for known vulnerabilities in container images using the configured scanner(s).

The scan runs asynchronously in the background. After initiating the scan, you can:
  - Check the scan progress with 'harbor-cli scan-all metrics'
  - View results through the Harbor UI 

Important considerations:
  - This operation can be resource intensive on large registries
  - Scanning many artifacts simultaneously may impact system performance
  - The time to complete depends on the number and size of artifacts
  - Only one scan-all operation can run at a time

Examples:
  # Start scanning all artifacts immediately
  harbor-cli scan-all run

  # Start scanning and monitor progress
  harbor-cli scan-all run && watch -n 0.2 harbor-cli scan-all metrics

The scan progress and results can be monitored through the metrics command
or through the Harbor web interface.`,
		Args: cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Info("Initiating manual scan of all artifacts")
			// Random cron expression and random time need to be passed to the API, even though they are not used, otherwise it returns bad request
			randomCron := "0 * * * * *"
			randomTime := strfmt.DateTime{}
			err := api.CreateScanAllSchedule(models.ScheduleObj{Type: "Manual", Cron: randomCron, NextScheduledTime: randomTime})
			if err != nil {
				return fmt.Errorf("failed to start scan all operation: %v", utils.ParseHarborErrorMsg(err))
			}
			logrus.Info("Successfully started scan all operation")
			return nil
		},
	}

	return cmd
}
