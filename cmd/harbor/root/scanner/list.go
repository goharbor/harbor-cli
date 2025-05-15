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
package scanner

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/scanner/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListScannerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List scanners",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			scannersResp, err := api.ListScanners()
			if err != nil {
				return fmt.Errorf("failed to list scanners: %v", err)
			}

			scanners := scannersResp.Payload
			if len(scanners) == 0 {
				log.Info("No scanners found")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(scanners, formatFlag)
				if err != nil {
					return err
				}
			} else {
				list.ListScanners(scanners)
			}

			return nil
		},
	}
	return cmd
}
