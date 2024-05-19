package project_tests

import (
	"context"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	cli "github.com/goharbor/harbor-cli/cmd/harbor/root/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunDeleteProjectSuccess(t *testing.T) {
	mockProject := new(MockProject)
	ctx := context.Background()
	projectName := "test"

	mockProject.On("DeleteProject", mock.Anything, mock.Anything).Return(&project.DeleteProjectOK{}, nil)

	err := cli.RunDeleteProject(projectName, ctx, mockProject)

	assert.NoError(t, err)
	mockProject.AssertExpectations(t)
}

func TestRunDeleteProjectFailed(t *testing.T) {
	mockProject := new(MockProject)
	ctx := context.Background()
	projectName := ""

	mockProject.On("DeleteProject", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	err := cli.RunDeleteProject(projectName, ctx, mockProject)

	assert.Error(t, err)
	mockProject.AssertExpectations(t)
}
