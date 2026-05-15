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
package instance

import (
	"fmt"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/instance/create"
	"github.com/spf13/cobra"
)

func CreateInstanceCommand() *cobra.Command {
	var opts create.CreateView
	var authUsername, authPassword, authToken string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new preheat provider instance in Harbor",
		Long: `Create a new preheat provider instance within Harbor for distributing container images.
The instance can be an external service such as Dragonfly, Kraken, or any custom provider.
You will need to provide the instance's name, vendor, endpoint, and optionally other details such as authentication and security options.`,
		Example: `  harbor-cli instance create --name my-instance --provider dragonfly --url http://dragonfly.local --description "My preheat provider instance" --enable=true`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var instanceName string

			if opts.Name != "" && opts.Vendor != "" && opts.Endpoint != "" {
				formattedEndpoint := utils.FormatUrl(opts.Endpoint)
				if err := utils.ValidateURL(formattedEndpoint); err != nil {
					return err
				}
				opts.Endpoint = formattedEndpoint

				opts.AuthMode = strings.ToUpper(strings.TrimSpace(opts.AuthMode))

				switch opts.AuthMode {
				case "BASIC":
					if authUsername == "" || authPassword == "" {
						return fmt.Errorf("username and password are required when authmode is BASIC. Use --auth-username and --auth-password flags")
					}
					opts.AuthInfo = map[string]string{
						"username": authUsername,
						"password": authPassword,
					}
				case "OAUTH":
					if authToken == "" {
						return fmt.Errorf("token is required when authmode is OAUTH. Use --auth-token flag")
					}
					opts.AuthInfo = map[string]string{
						"token": authToken,
					}
				case "NONE":
					// Auth credentials are ignored when authmode is NONE
				default:
					return fmt.Errorf("invalid authmode '%s'. Valid options: NONE, BASIC, OAUTH", opts.AuthMode)
				}

				err = api.CreateInstance(opts)
				instanceName = opts.Name
			} else {
				createView := &create.CreateView{
					Name:        opts.Name,
					Vendor:      opts.Vendor,
					Description: opts.Description,
					Endpoint:    opts.Endpoint,
					Insecure:    opts.Insecure,
					Enabled:     opts.Enabled,
					AuthMode:    opts.AuthMode,
				}
				err = createInstanceView(createView)
				instanceName = createView.Name
			}

			if err != nil {
				return fmt.Errorf("failed to create instance: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Printf("Instance '%s' created successfully\n", instanceName)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "Name of the instance")
	flags.StringVarP(&opts.Vendor, "provider", "p", "", "Provider for the instance (e.g. dragonfly, kraken)")
	flags.StringVarP(&opts.Endpoint, "url", "u", "", "Endpoint URL for the instance")
	flags.StringVarP(&opts.Description, "description", "d", "", "Description of the instance")
	flags.BoolVarP(&opts.Insecure, "insecure", "i", false, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.BoolVarP(&opts.Enabled, "enable", "", true, "Whether the instance is enabled or not")
	flags.StringVarP(&opts.AuthMode, "authmode", "a", "NONE", "Authentication mode (NONE, BASIC, OAUTH)")
	flags.StringVar(&authUsername, "auth-username", "", "Username for BASIC authentication")
	flags.StringVar(&authPassword, "auth-password", "", "Password for BASIC authentication")
	flags.StringVar(&authToken, "auth-token", "", "Token for OAUTH authentication")

	return cmd
}

func createInstanceView(createView *create.CreateView) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	if err := create.CreateInstanceView(createView); err != nil {
		return err
	}
	return api.CreateInstance(*createView)
}
