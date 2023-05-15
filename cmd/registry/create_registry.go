package registry

import (
	"context"

	"github.com/akshatdalton/harbor-cli/cmd/constants"
	"github.com/akshatdalton/harbor-cli/cmd/utils"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/spf13/cobra"
)

type createRegistrytOptions struct {
	name        string
	_type       string
	url         string
	description string
	insecure    bool
	credential  struct {
		accessKey    string
		accessSecret string
		_type        string
	}
}

// NewCreateRegistryCommand creates a new `harbor create registry` command
func NewCreateRegistryCommand() *cobra.Command {
	var opts createRegistrytOptions

	cmd := &cobra.Command{
		Use:   "registry",
		Short: "create registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return runCreateRegistry(opts, credentialName)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.name, "name", "", "", "Name of the registry")
	flags.StringVarP(&opts._type, "type", "", "harbor", "Type of the registry")
	flags.StringVarP(&opts.url, "url", "", "", "Registry endpoint URL")
	flags.StringVarP(&opts.description, "description", "", "", "Description of the registry")
	flags.BoolVarP(&opts.insecure, "insecure", "", true, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.StringVarP(&opts.credential.accessKey, "credential-access-key", "", "", "Access key, e.g. user name when credential type is 'basic'")
	flags.StringVarP(&opts.credential.accessKey, "credential-access-secret", "", "", "Access secret, e.g. password when credential type is 'basic'")
	flags.StringVarP(&opts.credential._type, "credential-type", "", "basic", "Credential type, such as 'basic', 'oauth'")

	return cmd
}

func runCreateRegistry(opts createRegistrytOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Registry.CreateRegistry(ctx, &registry.CreateRegistryParams{Registry: &models.Registry{Credential: &models.RegistryCredential{AccessKey: opts.credential.accessKey, AccessSecret: opts.credential.accessSecret, Type: opts.credential._type}, Description: opts.description, Insecure: opts.insecure, Name: opts.name, Type: opts._type, URL: opts.url}})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
