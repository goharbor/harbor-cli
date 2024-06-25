package registry

import (
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/registry/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewUpdateRegistryCommand creates a new `harbor update registry` command
func UpdateRegistryCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update registry",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var registryId int64

			updateView := &create.CreateView{
				Name:        opts.Name,
				Type:        opts.Type,
				Description: opts.Description,
				URL:         opts.URL,
				Credential: create.RegistryCredential{
					AccessKey:    opts.Credential.AccessKey,
					Type:         opts.Credential.Type,
					AccessSecret: opts.Credential.AccessSecret,
				},
				Insecure: opts.Insecure,
			}

			if len(args) > 0 {
				registryId, err = strconv.ParseInt(args[0], 10, 64)
			} else {
				registryId = prompt.GetRegistryNameFromUser()
			}

			if err != nil {
				log.Errorf("failed to parse registry id: %v", err)
			}

			if opts.Name != "" && opts.Type != "" && opts.URL != "" {
				err = api.UpdateRegistry(updateView, registryId)
			} else {
				err = updateRegistryView(updateView, registryId)
			}

			if err != nil {
				log.Errorf("failed to update registry: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the registry")
	flags.StringVarP(&opts.Type, "type", "", "", "Type of the registry")
	flags.StringVarP(&opts.URL, "url", "", "", "Registry endpoint URL")
	flags.StringVarP(&opts.Description, "description", "", "", "Description of the registry")
	flags.BoolVarP(&opts.Insecure, "insecure", "", true, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.StringVarP(&opts.Credential.AccessKey, "credential-access-key", "", "", "Access key, e.g. user name when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.AccessSecret, "credential-access-secret", "", "", "Access secret, e.g. password when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.Type, "credential-type", "", "", "Credential type, such as 'basic', 'oauth'")

	return cmd
}

func updateRegistryView(updateView *create.CreateView, projectID int64) error {
	if updateView == nil {
		updateView = &create.CreateView{}
	}

	create.CreateRegistryView(updateView)
	return api.UpdateRegistry(updateView, projectID)
}
