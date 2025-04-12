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
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/scanner/metadata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func MetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata [scanner-id]",
		Short: "Get scanner metadata by ID",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var registrationID string
			if len(args) > 0 {
				registrationID = args[0]
			} else {
				registrationID = prompt.GetScannerIdFromUser()
			}

			meta, err := api.GetScannerMetadata(registrationID)
			if err != nil {
				return fmt.Errorf("failed to get scanner metadata: %v", err)
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(meta, formatFlag)
				if err != nil {
					return err
				}
			} else {
				metadata.DisplayScannerMetadata(meta.Payload)
			}

			return nil
		},
	}

	return cmd
}
