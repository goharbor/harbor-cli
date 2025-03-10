// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/search"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func RepoDelete(projectNameOrID, repoName string, useProjectID bool) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	projectName := projectNameOrID
	if useProjectID {
		project, err := GetProject(projectNameOrID, useProjectID)
		if err != nil {
			return err
		}
		projectName = project.Payload.Name
	}
	_, err = client.Repository.DeleteRepository(ctx, &repository.DeleteRepositoryParams{ProjectName: projectName, RepositoryName: repoName})

	if err != nil {
		return err
	}

	log.Infof("Repository %s/%s deleted successfully", projectName, repoName)
	return nil
}

func RepoView(projectName, repoName string) (*repository.GetRepositoryOK, error) {
	ctx, client, err := utils.ContextWithClient()
	var response = &repository.GetRepositoryOK{}
	if err != nil {
		return response, err
	}

	response, err = client.Repository.GetRepository(ctx, &repository.GetRepositoryParams{ProjectName: projectName, RepositoryName: repoName})

	if err != nil {
		return response, err
	}

	return response, nil
}

func ListRepository(projectNameOrID string, useProjectID bool) (repository.ListRepositoriesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return repository.ListRepositoriesOK{}, err
	}
	projectName := projectNameOrID

	if useProjectID {
		project, err := GetProject(projectNameOrID, useProjectID)
		if err != nil {
			return repository.ListRepositoriesOK{}, err
		}
		projectName = project.Payload.Name
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
