package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type getProjectOptions struct {
	projectNameOrID string
}

// GetProjectCommand creates a new `harbor get project` command
func ViewCommand() *cobra.Command {
	var opts getProjectOptions

	cmd := &cobra.Command{
		Use:   "view [NAME|ID]",
		Short: "get project by name or id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.projectNameOrID = args[0]
			return runGetProject(opts)
		},
	}

	return cmd
}

func runGetProject(opts getProjectOptions) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.GetProject(ctx, &project.GetProjectParams{ProjectNameOrID: opts.projectNameOrID})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
