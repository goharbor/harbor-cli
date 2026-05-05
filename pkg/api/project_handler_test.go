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
	"errors"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/search"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestGetProject_Error(t *testing.T) {
	original := getProjectFunc
	defer func() { getProjectFunc = original }()

	getProjectFunc = func(projectNameOrID string, useProjectID bool) (*project.GetProjectOK, error) {
		return nil, errors.New("connection refused")
	}

	_, err := GetProject("test-project", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}

func TestGetProject_Success(t *testing.T) {
	original := getProjectFunc
	defer func() { getProjectFunc = original }()

	getProjectFunc = func(projectNameOrID string, useProjectID bool) (*project.GetProjectOK, error) {
		return &project.GetProjectOK{
			Payload: &models.Project{
				ProjectID: 42,
				Name:      "test-project",
			},
		}, nil
	}

	resp, err := GetProject("test-project", false)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-project", resp.Payload.Name)
	assert.Equal(t, int32(42), resp.Payload.ProjectID)
}

func TestGetProjectIDFromName_Error(t *testing.T) {
	original := getProjectFunc
	defer func() { getProjectFunc = original }()

	getProjectFunc = func(projectNameOrID string, useProjectID bool) (*project.GetProjectOK, error) {
		return nil, errors.New("project not found")
	}

	_, err := GetProjectIDFromName("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project not found")
}

func TestGetProjectIDFromName_Success(t *testing.T) {
	original := getProjectFunc
	defer func() { getProjectFunc = original }()

	getProjectFunc = func(projectNameOrID string, useProjectID bool) (*project.GetProjectOK, error) {
		return &project.GetProjectOK{
			Payload: &models.Project{
				ProjectID: 99,
				Name:      "my-project",
			},
		}, nil
	}

	id, err := GetProjectIDFromName("my-project")
	assert.NoError(t, err)
	assert.Equal(t, int64(99), id)
}

func TestListProject_Error(t *testing.T) {
	original := listProjectFunc
	defer func() { listProjectFunc = original }()

	listProjectFunc = func(opts ...ListFlags) (project.ListProjectsOK, error) {
		return project.ListProjectsOK{}, errors.New("server timeout")
	}

	_, err := ListProject()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server timeout")
}

func TestListProject_Success(t *testing.T) {
	original := listProjectFunc
	defer func() { listProjectFunc = original }()

	listProjectFunc = func(opts ...ListFlags) (project.ListProjectsOK, error) {
		return project.ListProjectsOK{
			Payload: []*models.Project{
				{ProjectID: 1, Name: "project-a"},
				{ProjectID: 2, Name: "project-b"},
			},
		}, nil
	}

	resp, err := ListProject()
	assert.NoError(t, err)
	assert.Len(t, resp.Payload, 2)
	assert.Equal(t, "project-a", resp.Payload[0].Name)
}

func TestListProject_Pagination(t *testing.T) {
	original := listProjectFunc
	defer func() { listProjectFunc = original }()

	var capturedPage, capturedPageSize int64
	listProjectFunc = func(opts ...ListFlags) (project.ListProjectsOK, error) {
		if len(opts) > 0 {
			capturedPage = opts[0].Page
			capturedPageSize = opts[0].PageSize
		}
		return project.ListProjectsOK{Payload: []*models.Project{}}, nil
	}

	_, err := ListProject(ListFlags{Page: 2, PageSize: 10})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), capturedPage)
	assert.Equal(t, int64(10), capturedPageSize)
}

func TestDeleteProject_Error(t *testing.T) {
	original := deleteProjectFunc
	defer func() { deleteProjectFunc = original }()

	deleteProjectFunc = func(projectNameOrID string, forceDelete bool, useProjectID bool) error {
		return errors.New("unauthorized")
	}

	err := DeleteProject("test-project", false, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestDeleteProject_Success(t *testing.T) {
	original := deleteProjectFunc
	defer func() { deleteProjectFunc = original }()

	deleteProjectFunc = func(projectNameOrID string, forceDelete bool, useProjectID bool) error {
		return nil
	}

	err := DeleteProject("test-project", false, false)
	assert.NoError(t, err)
}

func TestSearchProject_Error(t *testing.T) {
	original := searchProjectFunc
	defer func() { searchProjectFunc = original }()

	searchProjectFunc = func(query string) (search.SearchOK, error) {
		return search.SearchOK{}, errors.New("search service unavailable")
	}

	_, err := SearchProject("test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "search service unavailable")
}

func TestSearchProject_Success(t *testing.T) {
	original := searchProjectFunc
	defer func() { searchProjectFunc = original }()

	searchProjectFunc = func(query string) (search.SearchOK, error) {
		return search.SearchOK{}, nil
	}

	resp, err := SearchProject("found")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestLogsProject_Error(t *testing.T) {
	original := logsProjectFunc
	defer func() { logsProjectFunc = original }()

	logsProjectFunc = func(projectName string, opts ...ListFlags) (*project.GetLogExtsOK, error) {
		return nil, errors.New("logs not available")
	}

	_, err := LogsProject("test-project")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "logs not available")
}

func TestLogsProject_Success(t *testing.T) {
	original := logsProjectFunc
	defer func() { logsProjectFunc = original }()

	logsProjectFunc = func(projectName string, opts ...ListFlags) (*project.GetLogExtsOK, error) {
		return &project.GetLogExtsOK{}, nil
	}

	resp, err := LogsProject("test-project")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
