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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/registry/update"
	"github.com/spf13/cobra"
)

func UpdateRegistryCommand() *cobra.Command {
	opts := &models.Registry{
		Credential: &models.RegistryCredential{},
	}

	cmd := &cobra.Command{
		Use:   "update [registry_name]",
		Short: "update registry",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var registryId int64

			if len(args) > 0 {
				registryId, err = api.GetRegistryIdByName(args[0])
				if err != nil {
					return fmt.Errorf("failed to get registry id: %v", err)
				}
			} else {
				registryId = prompt.GetRegistryNameFromUser()
			}

			existingRegistry, err := api.GetRegistryResponse(registryId)
			if err != nil {
				return fmt.Errorf("failed to get registry with ID %d: %v", registryId, err)
			}
			if existingRegistry == nil {
				return fmt.Errorf("registry is not found")
			}

			updateView := &models.Registry{
				Name:        existingRegistry.Name,
				Type:        existingRegistry.Type,
				Description: existingRegistry.Description,
				URL:         existingRegistry.URL,
				Insecure:    existingRegistry.Insecure,
				Credential: &models.RegistryCredential{
					AccessKey:    existingRegistry.Credential.AccessKey,
					AccessSecret: existingRegistry.Credential.AccessSecret,
					Type:         existingRegistry.Credential.Type,
				},
			}

			flags := cmd.Flags()
			if flags.Changed("name") {
				updateView.Name = opts.Name
			}
			if flags.Changed("type") {
				updateView.Type = opts.Type
			}
			if flags.Changed("description") {
				updateView.Description = opts.Description
			}
			if flags.Changed("url") {
				formattedUrl := utils.FormatUrl(opts.URL)
				if err := utils.ValidateURL(formattedUrl); err != nil {
					return err
				}
				updateView.URL = formattedUrl
			}
			if flags.Changed("insecure") {
				updateView.Insecure = opts.Insecure
			}
			if flags.Changed("credential-access-key") {
				updateView.Credential.AccessKey = opts.Credential.AccessKey
			}
			if flags.Changed("credential-access-secret") {
				updateView.Credential.AccessSecret = opts.Credential.AccessSecret
			}
			if flags.Changed("credential-type") {
				updateView.Credential.Type = opts.Credential.Type
			}

			update.UpdateRegistryView(updateView)
			err = api.UpdateRegistry(updateView, registryId)
			if err != nil {
				return fmt.Errorf("failed to update registry: %v", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "Name of the registry")
	flags.StringVarP(&opts.Type, "type", "t", "", "Type of the registry")
	flags.StringVarP(&opts.URL, "url", "u", "", "Registry endpoint URL")
	flags.StringVarP(&opts.Description, "description", "d", "", "Description of the registry")
	flags.BoolVarP(&opts.Insecure, "insecure", "i", false, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.StringVarP(&opts.Credential.AccessKey, "credential-access-key", "k", "", "Access key, e.g. user name when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.AccessSecret, "credential-access-secret", "s", "", "Access secret, e.g. password when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.Type, "credential-type", "", "", "Credential type, such as 'basic', 'oauth'")

	return cmd
}
