package project

import (
	"context"

	"github.com/akshatdalton/harbor-cli/cmd/constants"
	"github.com/akshatdalton/harbor-cli/cmd/utils"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/spf13/cobra"
)

type deleteProjectOptions struct {
	projectNameOrID string
}

// NewDeleteProjectCommand creates a new `harbor delete project` command
func NewDeleteProjectCommand() *cobra.Command {
	var opts deleteProjectOptions

	cmd := &cobra.Command{
		Use:   "project [NAME|ID]",
		Short: "delete project by name or id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.projectNameOrID = args[0]
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return runDeleteProject(opts, credentialName)
		},
	}

	return cmd
}

func runDeleteProject(opts deleteProjectOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.DeleteProject(ctx, &project.DeleteProjectParams{ProjectNameOrID: opts.projectNameOrID})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
