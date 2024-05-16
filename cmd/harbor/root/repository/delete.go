package repository

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RepoDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a repository",
		Example: `  harbor repository delete [project_name]/[repository_name]`,
		Long:    `Delete a repository within a project in Harbor`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				projectName, repoName := utils.ParseProjectRepo(args[0])
				err = runRepoDelete(projectName, repoName)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				err = runRepoDelete(projectName, repoName)
			}
			if err != nil {
				log.Errorf("failed to delete repository: %v", err)
			}
		},
	}
	return cmd
}

func runRepoDelete(projectName, repoName string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	_, err := client.Repository.DeleteRepository(ctx, &repository.DeleteRepositoryParams{ProjectName: projectName, RepositoryName: repoName})

	if err != nil {
		return err
	}

	log.Infof("Repository %s/%s deleted successfully", projectName, repoName)
	return nil
}
