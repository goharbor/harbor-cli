package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/repository/list"
	log "github.com/sirupsen/logrus"
)

func RepoDelete(projectName, repoName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Repository.DeleteRepository(ctx, &repository.DeleteRepositoryParams{ProjectName: projectName, RepositoryName: repoName})

	if err != nil {
		return err
	}

	log.Infof("Repository %s/%s deleted successfully", projectName, repoName)
	return nil
}

func RepoInfo(projectName, repoName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.Repository.GetRepository(ctx, &repository.GetRepositoryParams{ProjectName: projectName, RepositoryName: repoName})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.Payload)
	return nil
}

func ListRepository(ProjectName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.Repository.ListRepositories(ctx, &repository.ListRepositoriesParams{ProjectName: ProjectName})

	if err != nil {
		return err
	}

	list.ListRepositories(response.Payload)
	return nil

}
