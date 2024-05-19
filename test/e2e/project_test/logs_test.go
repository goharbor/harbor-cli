package project_tests

import (
	"context"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cli "github.com/goharbor/harbor-cli/cmd/harbor/root/project"
)

func TestRunLogsProject(t *testing.T) {
	mockProject := new(MockProject)

	ctx := context.Background()

	projectName := "test"

	mockProject.On("GetLogs", mock.Anything, mock.Anything).Return(&project.GetLogsOK{}, nil)

	response, err := cli.RunLogsProject(projectName, ctx, mockProject)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	mockProject.AssertExpectations(t)
}

func TestRunLogsProjectError(t *testing.T) {
	mockProject := new(MockProject)

	ctx := context.Background()

	projectName := ""

	mockProject.On("GetLogs", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	response, err := cli.RunLogsProject(projectName, ctx, mockProject)

	assert.Error(t, err)
	assert.Nil(t, response)

	mockProject.AssertExpectations(t)

}
