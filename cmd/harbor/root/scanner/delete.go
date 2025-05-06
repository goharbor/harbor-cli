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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [scanner-name]",
		Short: "Delete a scanner registration",
		Long: `Delete a scanner registration from Harbor.

You can:
  - Provide the scanner name as an argument to delete it directly, or
  - Omit the argument to select a scanner interactively.

Note: Deleting a scanner will permanently remove its registration and associated metadata from the system.

Examples:
  # Delete a scanner by name
  harbor scanner delete trivy-scanner

  # Interactively choose a scanner to delete
  harbor scanner delete`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var registrationID string
			if len(args) > 0 {
				scanner, err := api.GetScannerByName(args[0])
				if err != nil {
					return fmt.Errorf("failed to retrieve scanner by name %q: %v", args[0], err)
				}
				registrationID = scanner.UUID
			} else {
				registrationID = prompt.GetScannerIdFromUser()
			}

			err = api.DeleteScanner(registrationID)
			if err != nil {
				return fmt.Errorf("failed to delete scanner: %v", err)
			}
			log.Infof("Scanner deleted successfully")
			return nil
		},
	}
	return cmd
}
