package registry

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// UpdateRegistryCommand creates a new `harbor update registry` command
func UpdateRegistryCommand() *cobra.Command {
	var opts config.UpdateRegistrytOptions

	cmd := &cobra.Command{
		Use:   "registry [ID]",
		Short: "update registry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid argument: %s. Expected an integer.\n", args[0])
				return err
			}
			opts.Id = int64(id)

			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return api.RunUpdateRegistry(opts, credentialName)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the registry")
	flags.StringVarP(&opts.Type, "type", "", "", "Type of the registry")
	flags.StringVarP(&opts.Url, "url", "", "", "Registry endpoint URL")
	flags.StringVarP(&opts.Description, "description", "", "", "Description of the registry")
	flags.BoolVarP(&opts.Insecure, "insecure", "", true, "Whether or not the certificate will be verified when Harbor tries to access the server")
	flags.StringVarP(&opts.Credential.AccessKey, "credential-access-key", "", "", "Access key, e.g. user name when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.AccessKey, "credential-access-secret", "", "", "Access secret, e.g. password when credential type is 'basic'")
	flags.StringVarP(&opts.Credential.Type, "credential-type", "", "", "Credential type, such as 'basic', 'oauth'")

	return cmd
}