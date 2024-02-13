package project

import (
	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// CreateProjectCommand creates a new `harbor create project` command
func CreateProjectCommand() *cobra.Command {
	var opts config.CreateProjectOptions

	cmd := &cobra.Command{
		Use:   "project [NAME]",
		Short: "create project",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectName = args[0]
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return api.RunCreateProject(opts, credentialName)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.Public, "public", "", false, "Project is public or private")
	flags.Int64VarP(&opts.RegistryID, "registry-id", "", 1, "ID of referenced registry when creating the proxy cache project")
	flags.Int64VarP(&opts.StorageLimit, "storage-limit", "", -1, "Storage quota of the project")

	return cmd
}
