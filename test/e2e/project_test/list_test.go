package project_tests

import (
	"context"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	cli "github.com/goharbor/harbor-cli/cmd/harbor/root/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// // Tests for success in the RunListProject function
func TestRunListProjectSuccess(t *testing.T) {
	mockProject := new(MockProject)

	ctx := context.Background()

	opts := cli.ListProjectOptions{
		Name:       "test",
		Owner:      "test",
		Page:       1,
		PageSize:   10,
		Public:     true,
		Q:          "test",
		Sort:       "test",
		WithDetail: true,
	}
	mockProject.On("ListProjects", mock.Anything, mock.Anything).Return(&project.ListProjectsOK{}, nil)

	resp, err := cli.RunListProject(opts, ctx, mockProject)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockProject.AssertExpectations(t)
}

// // Tests for error in the RunListProject function if Name field is missing in the ListProjectOptions struct
func TestRunListProjectsError(t *testing.T) {
	mockProject := new(MockProject)

	ctx := context.Background()

	opts := cli.ListProjectOptions{
		Owner:      "test",
		Page:       1,
		PageSize:   10,
		Public:     true,
		Q:          "test",
		Sort:       "test",
		WithDetail: true,
	}

	mockProject.On("ListProjects", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	resp, err := cli.RunListProject(opts, ctx, mockProject)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockProject.AssertExpectations(t)
}
