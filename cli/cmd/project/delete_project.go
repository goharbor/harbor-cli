package project

import (
	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// DeleteProjectCommand creates a new `harbor delete project` command
func DeleteProjectCommand() *cobra.Command {
	var opts config.DeleteProjectOptions

	cmd := &cobra.Command{
		Use:   "project [NAME|ID]",
		Short: "delete project by name or id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectNameOrID = args[0]
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return api.RunDeleteProject(opts, credentialName)
		},
	}

	return cmd
}
