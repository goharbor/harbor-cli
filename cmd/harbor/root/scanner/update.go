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

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/scanner/create"
	"github.com/goharbor/harbor-cli/pkg/views/scanner/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "update [scanner-name]",
		Short: "Update a scanner registration",
		Long: `Update the fields of an existing scanner registration.

You can pass the scanner name as an argument, or the CLI will prompt you to enter a scanner ID.
Only the fields passed through flags will be updated; other fields will retain their existing values.`,
		Example: `
  # Update description and URL for a scanner named 'trivy-scanner'
  harbor scanner update trivy-scanner --description "Updated scanner" --url "http://trivy.local:8080"

  # Change the authentication method and credential
  harbor scanner update trivy-scanner --auth Basic --credential "base64encodedAuth"

  # Disable the scanner and rename it
  harbor scanner update trivy-scanner --name "trivy-secure" --disabled

  # If no name is passed, you'll be prompted to enter a Scanner ID
  harbor scanner update --description "Updated via ID prompt"
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			resp, err := api.GetScanner(registrationID)
			if err != nil {
				return fmt.Errorf("scanner not found with ID %q: %v", registrationID, err)
			}
			existing := resp.GetPayload()

			updateView := &models.ScannerRegistration{
				Name:             existing.Name,
				Description:      existing.Description,
				Auth:             existing.Auth,
				AccessCredential: existing.AccessCredential,
				URL:              existing.URL,
				Disabled:         existing.Disabled,
				SkipCertVerify:   existing.SkipCertVerify,
				UseInternalAddr:  existing.UseInternalAddr,
			}

			flags := cmd.Flags()
			if flags.Changed("name") {
				updateView.Name = opts.Name
			}
			if flags.Changed("description") {
				updateView.Description = opts.Description
			}
			if flags.Changed("auth") {
				updateView.Auth = opts.Auth
			}
			if flags.Changed("credential") {
				updateView.AccessCredential = opts.AccessCredential
			}
			if flags.Changed("url") {
				formattedUrl := utils.FormatUrl(opts.URL)
				if err := utils.ValidateURL(formattedUrl); err != nil {
					return err
				}
				updateView.URL = strfmt.URI(formattedUrl)
			}
			if flags.Changed("disabled") {
				updateView.Disabled = &opts.Disabled
			}
			if flags.Changed("skip-cert-verification") {
				updateView.SkipCertVerify = &opts.SkipCertVerify
			}
			if flags.Changed("use-internal-addr") {
				updateView.UseInternalAddr = &opts.UseInternalAddr
			}

			update.UpdateScannerView(updateView)

			err = api.UpdateScanner(registrationID, *updateView)
			if err != nil {
				return fmt.Errorf("failed to update scanner: %v", err)
			}

			log.Infof("Scanner %q updated successfully", updateView.Name)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Name, "name", "", "New name for the scanner")
	flags.StringVar(&opts.Description, "description", "", "New description for the scanner")
	flags.StringVar(&opts.Auth, "auth", "", "Authentication method [None|Basic|Bearer|X-ScannerAdapter-API-Key]")
	flags.StringVar(&opts.AccessCredential, "credential", "", "Authorization header for the Scanner Adapter API")
	flags.StringVar(&opts.URL, "url", "", "Base URL of the scanner adapter")
	flags.BoolVar(&opts.Disabled, "disabled", false, "Disable the scanner registration")
	flags.BoolVar(&opts.SkipCertVerify, "skip-cert-verification", false, "Skip certificate verification in HTTP requests")
	flags.BoolVar(&opts.UseInternalAddr, "use-internal-addr", false, "Use internal registry address for scanning")

	return cmd
}
