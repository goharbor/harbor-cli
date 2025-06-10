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
