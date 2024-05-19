package project_tests

import (
	"context"
	"errors"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"

	"github.com/goharbor/harbor-cli/pkg/views/project/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cli "github.com/goharbor/harbor-cli/cmd/harbor/root/project"
)

func TestRunCreateProjectSuccess(t *testing.T) {
	// Create a mock for ProjectInterface
	mockProject := new(MockProject)

	ctx := context.Background()

	opts := create.CreateView{
		ProjectName:  "test",
		Public:       true,
		RegistryID:   "1",
		StorageLimit: "100",
	}
	mockProject.On("CreateProject", mock.Anything, mock.Anything).Return(&project.CreateProjectCreated{}, nil)
	err := cli.RunCreateProject(ctx, opts, mockProject)

	assert.NoError(t, err)
	mockProject.AssertExpectations(t)
}
func TestRunCreateProjectMissingProjectName(t *testing.T) {
	mockProject := new(MockProject)
	ctx := context.Background()

	opts := create.CreateView{
		ProjectName:  "",
		Public:       true,
		RegistryID:   "1",
		StorageLimit: "100",
	}

	mockProject.On("CreateProject", mock.Anything, mock.Anything).Return(nil, errors.New("project name is missing"))

	err := cli.RunCreateProject(ctx, opts, mockProject)

	assert.Error(t, err)
	assert.Equal(t, "project name is missing", err.Error())
	mockProject.AssertExpectations(t)
}
