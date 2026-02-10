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
package registry

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/registry/create"
	"github.com/spf13/cobra"
)

func CreateRegistryCommand() *cobra.Command {
	var opts api.CreateRegView

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create registry",
		Example: "harbor registry create",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			createView := &api.CreateRegView{
				Name:        opts.Name,
				Type:        opts.Type,
				Description: opts.Description,
				URL:         opts.URL,
				Credential: api.RegistryCredential{
					AccessKey:    opts.Credential.AccessKey,
					Type:         opts.Credential.Type,
					AccessSecret: opts.Credential.AccessSecret,
				},
				Insecure: opts.Insecure,
			}

			if opts.Name != "" && opts.Type != "" && opts.URL != "" {
				formattedUrl := utils.FormatUrl(opts.URL)
				err = utils.ValidateURL(formattedUrl)
				if err != nil {
					return err
				}
				opts.URL = formattedUrl
				err = api.CreateRegistry(opts)
			} else {
				err = createRegistryView(createView)
			}

			if err != nil {
				return fmt.Errorf("failed to create registry: %v", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the registry")
	flags.StringVarP(&opts.Type, "type", "", "", "Type of the registry")
	flags.StringVarP(&opts.URL, "url", "", "", "Registry endpoint URL")
	flags.StringVarP(&opts.Description, "description", "", "", "Description of the registry")
	flags.BoolVarP(
		&opts.Insecure,
		"insecure",
		"",
		true,
		"Whether Harbor will verify the server certificate",
	)
	flags.StringVarP(
		&opts.Credential.AccessKey,
		"credential-access-key",
		"",
		"",
		"Access key, e.g. user name when credential type is 'basic'",
	)
	flags.StringVarP(
		&opts.Credential.AccessSecret,
		"credential-access-secret",
		"",
		"",
		"Access secret, e.g. password when credential type is 'basic'",
	)
	flags.StringVarP(
		&opts.Credential.Type,
		"credential-type",
		"",
		"basic",
		"Credential type, such as 'basic', 'oauth'",
	)

	return cmd
}

func createRegistryView(createView *api.CreateRegView) error {
	if createView == nil {
		createView = &api.CreateRegView{}
	}

	create.CreateRegistryView(createView)
	return api.CreateRegistry(*createView)
}
