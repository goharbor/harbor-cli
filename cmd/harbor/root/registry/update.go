package registry

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/registry/update"
	log "github.com/sirupsen/logrus"
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
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var registryId int64

			if len(args) > 0 {
				registryId, err = api.GetRegistryIdByName(args[0])
				if err != nil {
					log.Errorf("failed to get registry id: %v", err)
					return
				}
			} else {
				registryId = prompt.GetRegistryNameFromUser()
			}

			existingRegistry := api.GetRegistryResponse(registryId)
			if existingRegistry == nil {
				log.Errorf("registry is not found")
				return
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
				updateView.URL = opts.URL
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
				log.Errorf("failed to update registry: %v", err)
				return
			}
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
	flags.StringVarP(&opts.Credential.Type, "credential-type", "c", "", "Credential type, such as 'basic', 'oauth'")

	return cmd
}
