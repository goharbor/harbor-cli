package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
)

type ProjectInterface interface {
	CreateProject(context.Context, *project.CreateProjectParams) (*project.CreateProjectCreated, error)
	DeleteProject(context.Context, *project.DeleteProjectParams) (*project.DeleteProjectOK, error)
	ListProjects(context.Context, *project.ListProjectsParams) (*project.ListProjectsOK, error)
	GetLogs(context.Context, *project.GetLogsParams) (*project.GetLogsOK, error)
}
