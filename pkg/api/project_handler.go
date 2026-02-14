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
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/search"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/create"
	log "github.com/sirupsen/logrus"
)

func CreateProject(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	registryID := new(int64)
	*registryID, _ = strconv.ParseInt(opts.RegistryID, 10, 64)

	if !opts.ProxyCache {
		registryID = nil
	}

	storageLimit, _ := strconv.ParseInt(opts.StorageLimit, 10, 64)

	public := strconv.FormatBool(opts.Public)

	response, err := client.Project.CreateProject(ctx, &project.CreateProjectParams{Project: &models.ProjectReq{ProjectName: opts.ProjectName, RegistryID: registryID, StorageLimit: &storageLimit, Public: &opts.Public, Metadata: &models.ProjectMetadata{Public: public}}})
	if err != nil {
		return err
	}

	if response != nil {
		log.Info("Project created successfully")
	}
	return nil
}

func GetProject(projectNameOrID string, useProjectID bool) (*project.GetProjectOK, error) {
	ctx, client, err := utils.ContextWithClient()
	response := &project.GetProjectOK{}

	if err != nil {
		return response, err
	}
	useResourceName := !useProjectID

	response, err = client.Project.GetProject(ctx, &project.GetProjectParams{
		ProjectNameOrID: projectNameOrID,
		XIsResourceName: &useResourceName,
	})
	if err != nil {
		return response, err
	}

	return response, nil
}

func GetProjectIDFromName(projectName string) (int64, error) {
	proj, err := GetProject(projectName, false)
	if err != nil {
		return 0, err
	}

	return int64(proj.Payload.ProjectID), nil
}

func DeleteProject(projectNameOrID string, forceDelete bool, useProjectID bool) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	if forceDelete {
		var resp repository.ListRepositoriesOK

		project, err := GetProject(projectNameOrID, useProjectID)
		if err != nil {
			log.Errorf("failed to get project name: %v", err)
			return err
		}
		projectName := project.Payload.Name

		immutables, err := ListImmutable(projectName)
		if err != nil {
			log.Errorf("failed to list immutables for project: %v", err)
			return err
		}
		for _, rule := range immutables.Payload {
			err = DeleteImmutable(projectName, rule.ID)
			if err != nil {
				log.Errorf("failed to delete tag immutable rule: %v", err)
				return err
			}
		}

		resp, err = ListRepository(projectNameOrID, useProjectID)
		if err != nil {
			log.Errorf("failed to list repositories: %v", err)
			return err
		}

		for _, repo := range resp.Payload {
			_, repoName, err := utils.ParseProjectRepo(repo.Name)
			if err != nil {
				log.Errorf("failed to parse project/repo: %v", err)
				return err
			}
			err = RepoDelete(projectNameOrID, repoName, useProjectID)
			if err != nil {
				log.Errorf("failed to delete repository: %v", err)
				return err
			}
		}
	}
	useProjectName := !useProjectID
	_, err = client.Project.DeleteProject(ctx, &project.DeleteProjectParams{
		ProjectNameOrID: projectNameOrID,
		XIsResourceName: &useProjectName,
	})
	if err != nil {
		return err
	}

	log.Infof("Project %s deleted successfully", projectNameOrID)
	return nil
}

func ListProject(opts ...ListFlags) (project.ListProjectsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return project.ListProjectsOK{}, err
	}
	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}
	response, err := client.Project.ListProjects(ctx, &project.ListProjectsParams{Page: &listFlags.Page, PageSize: &listFlags.PageSize, Q: &listFlags.Q, Sort: &listFlags.Sort, Name: &listFlags.Name, Public: &listFlags.Public})
	if err != nil {
		return project.ListProjectsOK{}, err
	}
	return *response, nil
}

func ListAllProjects(opts ...ListFlags) (project.ListProjectsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return project.ListProjectsOK{}, err
	}
	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}
	response, err := client.Project.ListProjects(ctx, &project.ListProjectsParams{Page: &listFlags.Page, PageSize: &listFlags.PageSize, Q: &listFlags.Q, Sort: &listFlags.Sort, Name: &listFlags.Name})
	if err != nil {
		return project.ListProjectsOK{}, err
	}
	return *response, nil
}

func SearchProject(query string) (search.SearchOK, error) {
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

func LogsProject(projectName string, opts ...ListFlags) (*project.GetLogExtsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Project.GetLogExts(ctx, &project.GetLogExtsParams{
		ProjectName: projectName,
		Page:        &listFlags.Page,
		PageSize:    &listFlags.PageSize,
		Q:           &listFlags.Q,
		Sort:        &listFlags.Sort,
		Context:     ctx,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func CheckProject(projectName string) (bool, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return false, err
	}

	response, err := client.Project.HeadProject(ctx, &project.HeadProjectParams{
		ProjectName: projectName,
		Context:     ctx,
	})
	if err != nil {
		if utils.ParseHarborErrorCode(err) == "404" {
			return false, nil
		}
		return false, err
	}

	return response.IsSuccess(), nil
}
