package project_tests

import (
	"context"
	"errors"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/stretchr/testify/mock"
)

type MockProject struct {
	mock.Mock
}

func (m *MockProject) CreateProject(ctx context.Context, params *project.CreateProjectParams) (*project.CreateProjectCreated, error) {
	args := m.Called(ctx, params)
	if params.Project.Public == nil {
		return nil, errors.New("public field is missing")
	}
	return &project.CreateProjectCreated{}, args.Error(1)
}
func (m *MockProject) DeleteProject(ctx context.Context, params *project.DeleteProjectParams) (*project.DeleteProjectOK, error) {
	args := m.Called(ctx, params)
	return &project.DeleteProjectOK{}, args.Error(1)
}

func (m *MockProject) ListProjects(ctx context.Context, params *project.ListProjectsParams) (*project.ListProjectsOK, error) {
	args := m.Called(ctx, params)
	return &project.ListProjectsOK{}, args.Error(1)
}

func (m *MockProject) GetLogs(ctx context.Context, params *project.GetLogsParams) (*project.GetLogsOK, error) {
	args := m.Called(ctx, params)
	return &project.GetLogsOK{}, args.Error(1)
}
