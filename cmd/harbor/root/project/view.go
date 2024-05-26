package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type getProjectOptions struct {
	isID bool
}

// GetProjectCommand creates a new `harbor get project` command
func ViewCommand() *cobra.Command {
	var opts getProjectOptions

	cmd := &cobra.Command{
		Use:   "view [NAME|ID]",
		Short: "get project by name or id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetProject(args[0], opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.isID, "id", false, "Get project by id")

	return cmd
}

func runGetProject(projectNameOrID string, opts getProjectOptions) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	isName := !opts.isID
	response, err := client.Project.GetProject(ctx, &project.GetProjectParams{XIsResourceName: &isName, ProjectNameOrID: projectNameOrID})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
