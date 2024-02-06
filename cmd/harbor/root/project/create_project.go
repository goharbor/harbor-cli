package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type createProjectOptions struct {
	projectName  string
	public       bool
	registryID   int64
	storageLimit int64
}

// CreateProjectCommand creates a new `harbor create project` command
func CreateProjectCommand() *cobra.Command {
	var opts createProjectOptions

	cmd := &cobra.Command{
		Use:   "project",
		Short: "create project",
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return runCreateProject(opts, credentialName)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.projectName, "name", "", "", "Name of the project")
	flags.BoolVarP(&opts.public, "public", "", true, "Project is public or private")
	flags.Int64VarP(&opts.registryID, "registry-id", "", 1, "ID of referenced registry when creating the proxy cache project")
	flags.Int64VarP(&opts.storageLimit, "storage-limit", "", -1, "Storage quota of the project")

	return cmd
}

func runCreateProject(opts createProjectOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.CreateProject(ctx, &project.CreateProjectParams{Project: &models.ProjectReq{ProjectName: opts.projectName, Public: &opts.public, RegistryID: &opts.registryID, StorageLimit: &opts.storageLimit}})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
