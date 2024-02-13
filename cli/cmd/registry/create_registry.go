package registry

import (
	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// CreateRegistryCommand creates a new `harbor create registry` command
func CreateRegistryCommand() *cobra.Command {
	var opts config.CreateRegistrytOptions

	cmd := &cobra.Command{
		Use:   "registry",
		Short: "create registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return api.RunCreateRegistry(opts, credentialName)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "Name of the registry")
	flags.StringVarP(&opts.Type, "type", "", "harbor", "Type of the registry")
	flags.StringVarP(&opts.Url, "url", "", "", "Registry endpoint URL")
	flags.StringVarP(&opts.Description, "description", "", "", "Description of the registry")
	flags.BoolVarP(&opts.Insecure, "insecure", "", true, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.StringVarP(&opts.Credential.AccessKey, "credential-access-key", "", "", "Access key, e.g. user name when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.AccessSecret, "credential-access-secret", "", "", "Access secret, e.g. password when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.Type, "credential-type", "", "basic", "Credential type, such as 'basic', 'oauth'")

	return cmd
}
