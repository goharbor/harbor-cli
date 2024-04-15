package registry

import (
	"context"
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type updateRegistrytOptions struct {
	id          int64
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

// NewUpdateRegistryCommand creates a new `harbor update registry` command
func UpdateRegistryCommand() *cobra.Command {
	var opts updateRegistrytOptions

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update registry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid argument: %s. Expected an integer.\n", args[0])
				return err
			}
			opts.id = int64(id)

			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return runUpdateRegistry(opts, credentialName)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.name, "name", "", "", "Name of the registry")
	flags.StringVarP(&opts._type, "type", "", "", "Type of the registry")
	flags.StringVarP(&opts.url, "url", "", "", "Registry endpoint URL")
	flags.StringVarP(&opts.description, "description", "", "", "Description of the registry")
	flags.BoolVarP(&opts.insecure, "insecure", "", true, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.StringVarP(&opts.credential.accessKey, "credential-access-key", "", "", "Access key, e.g. user name when credential type is 'basic'")
	flags.StringVarP(&opts.credential.accessKey, "credential-access-secret", "", "", "Access secret, e.g. password when credential type is 'basic'")
	flags.StringVarP(&opts.credential._type, "credential-type", "", "", "Credential type, such as 'basic', 'oauth'")

	return cmd
}

func runUpdateRegistry(opts updateRegistrytOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	registryUpdate := &models.RegistryUpdate{}

	if opts.credential.accessKey != "" {
		registryUpdate.AccessKey = &opts.credential.accessKey
	}

	if opts.credential.accessSecret != "" {
		registryUpdate.AccessSecret = &opts.credential.accessSecret
	}

	if opts.credential._type != "" {
		registryUpdate.CredentialType = &opts.credential._type
	}

	if opts.description != "" {
		registryUpdate.Description = &opts.description
	}

	if opts.name != "" {
		registryUpdate.Name = &opts.name
	}

	if opts.url != "" {
		registryUpdate.URL = &opts.url
	}

	registryUpdate.Insecure = &opts.insecure

	response, err := client.Registry.UpdateRegistry(ctx, &registry.UpdateRegistryParams{ID: opts.id, Registry: registryUpdate})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
