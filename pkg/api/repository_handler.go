package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/search"
	"github.com/goharbor/harbor-cli/pkg/utils"
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

func ListRepository(projectName string) (repository.ListRepositoriesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return repository.ListRepositoriesOK{}, err
	}

	response, err := client.Repository.ListRepositories(ctx, &repository.ListRepositoriesParams{ProjectName: projectName})

	if err != nil {
		return repository.ListRepositoriesOK{}, err
	}

	return *response, nil

}

func SearchRepository(query string) (search.SearchOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return search.SearchOK{}, err
	}

	response, err := client.Search.Search(ctx, &search.SearchParams{Q: query})
	if err != nil {
		return search.SearchOK{}, err
	}

	return *response, nil
}
