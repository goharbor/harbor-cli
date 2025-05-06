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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/scanner/view"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ViewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "view [scanner-name]",
		Short: "Display detailed information about a scanner registration",
		Long: `Display full details of a scanner registration in Harbor.

You can:
  - Provide a scanner name to view its details directly.
  - Omit the argument to select a scanner interactively by ID.

Supports custom output formats using the --output-format flag (e.g., json, yaml, table).

Examples:
  # View a specific scanner by name
  harbor scanner view trivy-scanner

  # Interactively choose a scanner to view
  harbor scanner view

  # View scanner in JSON format
  harbor scanner view trivy-scanner --output-format=json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var scanner *models.ScannerRegistration
			if len(args) > 0 {
				resp, err := api.GetScannerByName(args[0])
				if err != nil {
					return fmt.Errorf("failed to get scanner by name %q: %v", args[0], err)
				}
				scanner = &resp
			} else {
				id := prompt.GetScannerIdFromUser()
				resp, err := api.GetScanner(id)
				if err != nil {
					return fmt.Errorf("failed to get scanner by ID %q: %v", id, err)
				}
				scanner = resp.GetPayload()
			}

			outputFormat := viper.GetString("output-format")
			if outputFormat != "" {
				if err := utils.PrintFormat(scanner, outputFormat); err != nil {
					return fmt.Errorf("failed to format output: %v", err)
				}
			} else {
				view.ViewScanner(scanner)
			}
			return nil
		},
	}
}
