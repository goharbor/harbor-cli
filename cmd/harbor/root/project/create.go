package project

import (
	"context"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CreateProjectCommand creates a new `harbor create project` command
func CreateProjectCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "project",
		Short: "create project",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				log.Fatalf("failed to get credential name: %v", err)
			}

			createView := &create.CreateView{
				ProjectName:  opts.ProjectName,
				Public:       opts.Public,
				RegistryID:   opts.RegistryID,
				StorageLimit: opts.StorageLimit,
			}

			if opts.ProjectName != "" && opts.RegistryID != "" && opts.StorageLimit != "" {
				err = runCreateProject(opts, credentialName)
			} else {
				err = createProjectView(createView, credentialName)
			}

			if err != nil {
				log.Errorf("failed to create project: %v", err)
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.ProjectName, "name", "", "", "Name of the project")
	flags.BoolVarP(&opts.Public, "public", "", true, "Project is public or private")
	flags.StringVarP(&opts.RegistryID, "registry-id", "", "", "ID of referenced registry when creating the proxy cache project")
	flags.StringVarP(&opts.StorageLimit, "storage-limit", "", "-1", "Storage quota of the project")

	return cmd
}

func createProjectView(createView *create.CreateView, credentialName string) error {
	if createView == nil {
		createView = &create.CreateView{
			ProjectName:  "",
			Public:       true,
			RegistryID:   "",
			StorageLimit: "-1",
		}
	}

	create.CreateProjectView(createView)

	return runCreateProject(*createView, credentialName)

}

func runCreateProject(opts create.CreateView, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	registryID, _ := strconv.ParseInt(opts.RegistryID, 10, 64)

	storageLimit, _ := strconv.ParseInt(opts.StorageLimit, 10, 64)

	response, err := client.Project.CreateProject(ctx, &project.CreateProjectParams{Project: &models.ProjectReq{ProjectName: opts.ProjectName, Public: &opts.Public, RegistryID: &registryID, StorageLimit: &storageLimit}})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
