package api

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/search"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func RepoDelete(projectName, repoName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client context")
	}

	_, err = client.Repository.DeleteRepository(ctx, &repository.DeleteRepositoryParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
	})

	if err != nil {
		switch err.(type) {
		case *repository.DeleteRepositoryNotFound:
			return fmt.Errorf("repository not found: %s/%s", projectName, repoName)
		case *repository.DeleteRepositoryBadRequest:
			return fmt.Errorf("bad request while deleting repository: %s/%s", projectName, repoName)
		case *repository.DeleteRepositoryForbidden:
			return fmt.Errorf("forbidden to delete repository: %s/%s", projectName, repoName)
		case *repository.DeleteRepositoryInternalServerError:
			return fmt.Errorf("internal server error occurred while deleting repository: %s/%s", projectName, repoName)
		default:
			return fmt.Errorf("unknown error occurred while deleting repository: %v", err)
		}
	}

	log.Infof("Repository %s/%s deleted successfully", projectName, repoName)
	return nil
}

func RepoView(projectName, repoName string) (*repository.GetRepositoryOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client context")
	}

	response, err := client.Repository.GetRepository(ctx, &repository.GetRepositoryParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
	})

	if err != nil {
		switch err.(type) {
		case *repository.GetRepositoryNotFound:
			return nil, fmt.Errorf("repository not found: %s/%s", projectName, repoName)
		case *repository.GetRepositoryBadRequest:
			return nil, fmt.Errorf("bad request while getting repository: %s/%s", projectName, repoName)
		case *repository.GetRepositoryForbidden:
			return nil, fmt.Errorf("forbidden to get repository: %s/%s", projectName, repoName)
		case *repository.GetRepositoryInternalServerError:
			return nil, fmt.Errorf("internal server error occurred while getting repository: %s/%s", projectName, repoName)
		default:
			return nil, fmt.Errorf("unknown error occurred while getting repository: %v", err)
		}
	}

	log.Infof("Repository %s/%s details retrieved successfully", projectName, repoName)
	return response, nil
}

func ListRepository(projectName string) (repository.ListRepositoriesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return repository.ListRepositoriesOK{}, fmt.Errorf("failed to initialize client context")
	}

	response, err := client.Repository.ListRepositories(ctx, &repository.ListRepositoriesParams{
		ProjectName: projectName,
	})

	if err != nil {
		switch err.(type) {
		case *repository.ListRepositoriesNotFound:
			return repository.ListRepositoriesOK{}, fmt.Errorf("project not found: %s", projectName)
		case *repository.ListRepositoriesBadRequest:
			return repository.ListRepositoriesOK{}, fmt.Errorf("bad request while listing repositories: %s", projectName)
		case *repository.ListRepositoriesForbidden:
			return repository.ListRepositoriesOK{}, fmt.Errorf("forbidden to list repositories: %s", projectName)
		case *repository.ListRepositoriesInternalServerError:
			return repository.ListRepositoriesOK{}, fmt.Errorf("internal server error occurred while listing repositories: %s", projectName)
		default:
			return repository.ListRepositoriesOK{}, fmt.Errorf("unknown error occurred while listing repositories: %v", err)
		}
	}

	log.Infof("Repositories for project %s listed successfully", projectName)
	return *response, nil
}
func SearchRepository(query string) (search.SearchOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return search.SearchOK{}, fmt.Errorf("failed to initialize client context")
	}

	response, err := client.Search.Search(ctx, &search.SearchParams{Q: query})
	if err != nil {
		switch err.(type) {
		case *search.SearchInternalServerError:
			return search.SearchOK{}, fmt.Errorf("internal server error occurred while searching repositories: %s", query)
		default:
			return search.SearchOK{}, fmt.Errorf("unknown error occurred while searching repositories: %v", err)
		}
	}

	log.Infof("Repositories matching query '%s' retrieved successfully", query)
	return *response, nil
}
