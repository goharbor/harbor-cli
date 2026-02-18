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
	"github.com/goharbor/harbor-cli/pkg/views/scanner/create"
	"github.com/spf13/cobra"
)

func CreateScannerCommand() *cobra.Command {
	var opts create.CreateView
	var ping bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a scanner",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Name == "" || opts.Auth == "" || opts.URL == "" {
				create.CreateScannerView(&opts)
			} else {
				// Validate URL when provided via flags
				formattedUrl := utils.FormatUrl(opts.URL)
				if err := utils.ValidateURL(formattedUrl); err != nil {
					return err
				}
				opts.URL = formattedUrl
			}

			if ping {
				err := api.PingScanner(opts)
				if err != nil {
					return fmt.Errorf("failed to ping the scanner adapter: %v", err)
				}
			} else {
				err := api.CreateScanner(opts)
				if err != nil {
					return fmt.Errorf("failed to create scanner: %v", err.Error())
				}
			}
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
	flags.BoolVarP(&ping, "ping", "", false, "Ping the scanner adapter without creating it.")

	return cmd
}
