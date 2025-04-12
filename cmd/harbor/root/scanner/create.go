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
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/scanner/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateScannerCommand() *cobra.Command {
	var opts create.CreateView
	var ping bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a scanner",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if opts.Name == "" || opts.Description == "" || opts.Auth == "" || opts.AccessCredential == "" || opts.URL == "" {
				create.CreateScannerView(&opts)
			}

			if ping {
				err := api.PingScanner(opts)
				if err != nil {
					log.Errorf("failed to ping the scanner adapter: %v", err)
				}
			} else {
				err := api.CreateScanner(opts)
				if err != nil {
					log.Errorf("failed to create scanner: %v", err.Error())
				}
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the scanner")
	flags.StringVarP(&opts.Description, "des", "", "", "Description of the scanner")
	flags.StringVarP(&opts.Auth, "auth", "", "", "Authentication approach of the scanner [None|Basic|Bearer|X-ScannerAdapter-API-Key]")
	flags.StringVarP(&opts.AccessCredential, "cred", "", "", "HTTP Authorization header value sent with each request to the Scanner Adapter API")
	flags.StringVarP(&opts.URL, "url", "", "", "Base URL of the scanner adapter")
	flags.BoolVarP(&opts.Disabled, "disable", "", false, "Indicate whether the registration is enabled or not")
	flags.BoolVarP(&opts.SkipCertVerify, "skip", "", false, "Indicate if skip the certificate verification when sending HTTP requests")
	flags.BoolVarP(&opts.UseInternalAddr, "internal", "", false, "Indicate whether use internal registry addr for the scanner to pull content or not")
	flags.BoolVarP(&ping, "ping", "", false, "Ping the scanner adapter without creating it.")

	return cmd
}
