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
	"github.com/spf13/cobra"
)

func SetDefaultCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-default",
		Short:   "Set the default scanner for Harbor",
		Long:    `Set the scanner that will be used as the default in Harbor. This scanner will be used for all default scanning tasks unless another scanner is specified.`,
		Aliases: []string{"sd"},
		Example: `harbor scanner set-default [scanner-name]
		OR 
		harbor scanner set-default --id <scanner-id>`,
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
			err = api.SetDefaultScanner(registrationID)
			if err != nil {
				return fmt.Errorf("failed to set default scanner: %v", err)
			}
			fmt.Printf("Scanner %q successfully set as the default.\n", registrationID)
			return nil
		},
	}
	return cmd
}
