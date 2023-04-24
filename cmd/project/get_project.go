package project

import (
	"context"

	"github.com/akshatdalton/harbor-cli/cmd/constants"
	"github.com/akshatdalton/harbor-cli/cmd/utils"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/spf13/cobra"
)

type getProjectOptions struct {
	projectNameOrID string
}

// NewGetProjectCommand creates a new `harbor get project` command
func NewGetProjectCommand() *cobra.Command {
	var opts getProjectOptions

	cmd := &cobra.Command{
		Use:   "project [NAME|ID]",
		Short: "get project by name or id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.projectNameOrID = args[0]
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return runGetProject(opts, credentialName)
		},
	}

	return cmd
}

func runGetProject(opts getProjectOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.GetProject(ctx, &project.GetProjectParams{ProjectNameOrID: opts.projectNameOrID})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
