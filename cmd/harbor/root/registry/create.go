package registry

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/registry/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCreateRegistryCommand creates a new `harbor create registry` command
func CreateRegistryCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create registry",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			createView := &create.CreateView{
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

			if opts.Name != "" && opts.Type != "" && opts.URL != "" {
				err = api.CreateRegistry(opts)
			} else {
				err = createRegistryView(createView)
			}

			if err != nil {
				log.Errorf("failed to create registry: %v", err)
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the registry")
	flags.StringVarP(&opts.Type, "type", "", "harbor", "Type of the registry")
	flags.StringVarP(&opts.URL, "url", "", "", "Registry endpoint URL")
	flags.StringVarP(&opts.Description, "description", "", "", "Description of the registry")
	flags.BoolVarP(&opts.Insecure, "insecure", "", true, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.StringVarP(&opts.Credential.AccessKey, "credential-access-key", "", "", "Access key, e.g. user name when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.AccessSecret, "credential-access-secret", "", "", "Access secret, e.g. password when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.Type, "credential-type", "", "basic", "Credential type, such as 'basic', 'oauth'")

	return cmd
}

func createRegistryView(createView *create.CreateView) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateRegistryView(createView)
	return api.CreateRegistry(*createView)
}
