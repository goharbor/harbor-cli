package repository

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RepoInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "Get repository information",
		Example: `  harbor repo info <project_name>/<repo_name>`,
		Long:    `Get information of a particular repository in a project`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				projectName, repoName := utils.ParseProjectRepo(args[0])
				err = runRepoInfo(projectName, repoName)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				err = runRepoInfo(projectName, repoName)
			}
			if err != nil {
				log.Errorf("failed to get repository information: %v", err)
			}

		},
	}

	return cmd
}

func runRepoInfo(projectName, repoName string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	response, err := client.Repository.GetRepository(ctx, &repository.GetRepositoryParams{ProjectName: projectName, RepositoryName: repoName})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.Payload)
	return nil
}
