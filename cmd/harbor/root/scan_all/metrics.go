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
	view "github.com/goharbor/harbor-cli/pkg/views/scan-all/metrics"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetScanAllMetricsCommand() *cobra.Command {
	var scheduled bool

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Get the metrics of the latest scan all process",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Info("Retrieving scan all metrics")
			metrics, err := api.GetScanAllMetrics(scheduled)
			if err != nil {
				logrus.Errorf("Failed to retrieve scan all metrics: %v", utils.ParseHarborErrorMsg(err))
				return err
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(metrics, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				view.ViewScanMetrics(metrics)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	// latest scheduled metrics is deprecated in the API
	flags.BoolVarP(&scheduled, "scheduled", "s", false, "Get the metrics of the latest scheduled scan all process")

	return cmd
}
