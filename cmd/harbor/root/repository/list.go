package repository

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/repository/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListRepositoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list repositories within a project",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				err = runListRepository(args[0])
			} else {
				projectName := utils.GetProjectNameFromUser()
				err = runListRepository(projectName)
			}
			if err != nil {
				log.Errorf("failed to list repositories: %v", err)
			}
		},
	}

	return cmd
}

func runListRepository(ProjectName string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	response, err := client.Repository.ListRepositories(ctx, &repository.ListRepositoriesParams{ProjectName: ProjectName})

	if err != nil {
		return err
	}

	list.ListRepositories(response.Payload)
	return nil

}
