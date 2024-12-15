package api

import (
	"fmt"
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
		return fmt.Errorf("Failed to initialize client context for creating project")
	}
	registryID := new(int64)
	*registryID, _ = strconv.ParseInt(opts.RegistryID, 10, 64)

	if !opts.ProxyCache {
		registryID = nil
	}

	storageLimit, _ := strconv.ParseInt(opts.StorageLimit, 10, 64)
	public := strconv.FormatBool(opts.Public)

	response, err := client.Project.CreateProject(ctx, &project.CreateProjectParams{Project: &models.ProjectReq{
		ProjectName:  opts.ProjectName,
		RegistryID:   registryID,
		StorageLimit: &storageLimit,
		Public:       &opts.Public,
		Metadata:     &models.ProjectMetadata{Public: public},
	}})

	if err != nil {
		switch err.(type) {
		case *project.CreateProjectBadRequest:
			return fmt.Errorf("Bad request for creating project: %s", opts.ProjectName)
		case *project.CreateProjectConflict:
			return fmt.Errorf("Project %s already exists", opts.ProjectName)
		case *project.CreateProjectInternalServerError:
			return fmt.Errorf("Internal server error occurred while creating project: %s", opts.ProjectName)
		case *project.CreateProjectUnauthorized:
			return fmt.Errorf("Unauthorized to create project: %s", opts.ProjectName)
		default:
			return fmt.Errorf("Unknown error occurred while creating project %s: %v", opts.ProjectName, err)
		}
	}
	if response != nil {
		log.Info("Project created successfully")
	}
	return nil
}

func GetProject(projectName string) (*project.GetProjectOK, error) {
	ctx, client, err := utils.ContextWithClient()
	var response = &project.GetProjectOK{}
	if err != nil {
		return response, fmt.Errorf("Failed to initialize client context for getting project %s", projectName)
	}

	response, err = client.Project.GetProject(ctx, &project.GetProjectParams{ProjectNameOrID: projectName})
	if err != nil {
		switch err.(type) {
		case *project.GetProjectInternalServerError:
			return response, fmt.Errorf("Internal server error occurred while getting project %s", projectName)
		case *project.GetProjectDeletableNotFound:
			return response, fmt.Errorf("Project %s not found", projectName)
		case *project.GetProjectUnauthorized:
			return response, fmt.Errorf("Unauthorized to get project %s", projectName)
		default:
			return response, fmt.Errorf("Unknown error occurred while getting project %s: %v", projectName, err)
		}
	}

	return response, nil
}

func DeleteProject(projectName string, forceDelete bool) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("Failed to initialize client context for deleting project %s", projectName)
	}

	if forceDelete {
		var resp repository.ListRepositoriesOK
		resp, err = ListRepository(projectName)
		if err != nil {
			return fmt.Errorf("Failed to list repositories for project %s", projectName)
		}

		for _, repo := range resp.Payload {
			_, repoName := utils.ParseProjectRepo(repo.Name)
			err = RepoDelete(projectName, repoName)
			if err != nil {
				return fmt.Errorf("Failed to delete repository %s from project %s", repoName, projectName)
			}
		}
	}

	_, err = client.Project.DeleteProject(ctx, &project.DeleteProjectParams{ProjectNameOrID: projectName})
	if err != nil {
		switch err.(type) {
		case *project.DeleteProjectBadRequest:
			return fmt.Errorf("Invalid request to delete project %s", projectName)
		case *project.DeleteProjectNotFound:
			return fmt.Errorf("Project %s not found or already deleted", projectName)
		case *project.DeleteProjectForbidden:
			return fmt.Errorf("Insufficient permissions to delete project %s", projectName)
		case *project.DeleteProjectInternalServerError:
			return fmt.Errorf("Internal server error occurred while deleting project %s", projectName)
		default:
			return fmt.Errorf("Unknown error occurred while deleting project %s: %v", projectName, err)
		}
	}

	log.Infof("Project %s deleted successfully", projectName)
	return nil
}

func ListProject(opts ...ListFlags) (project.ListProjectsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return project.ListProjectsOK{}, fmt.Errorf("Failed to initialize client context for listing projects")
	}
	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}
	response, err := client.Project.ListProjects(ctx, &project.ListProjectsParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Q:        &listFlags.Q,
		Sort:     &listFlags.Sort,
		Name:     &listFlags.Name,
		Public:   &listFlags.Public,
	})
	if err != nil {
		switch err.(type) {
		case *project.ListProjectsUnauthorized:
			return project.ListProjectsOK{}, fmt.Errorf("Unauthorized access to list the projects")
		case *project.ListProjectsInternalServerError:
			return project.ListProjectsOK{}, fmt.Errorf("Internal server error occurred while listing projects")
		default:
			return project.ListProjectsOK{}, fmt.Errorf("Unknown error occurred while listing projects: %v", err)
		}
	}
	return *response, nil
}

func ListAllProjects(opts ...ListFlags) (project.ListProjectsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return project.ListProjectsOK{}, fmt.Errorf("Failed to initialize client context for listing all projects")
	}
	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}
	response, err := client.Project.ListProjects(ctx, &project.ListProjectsParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Q:        &listFlags.Q,
		Sort:     &listFlags.Sort,
		Name:     &listFlags.Name,
	})
	if err != nil {
		switch err.(type) {
		case *project.ListProjectsUnauthorized:
			return project.ListProjectsOK{}, fmt.Errorf("Unauthorized access to list the projects")
		case *project.ListProjectsInternalServerError:
			return project.ListProjectsOK{}, fmt.Errorf("Internal server error occurred while listing all projects")
		default:
			return project.ListProjectsOK{}, fmt.Errorf("Unknown error occurred while listing all projects: %v", err)
		}
	}
	return *response, nil
}

func SearchProject(query string) (search.SearchOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return search.SearchOK{}, fmt.Errorf("Failed to initialize client context for searching projects")
	}

	response, err := client.Search.Search(ctx, &search.SearchParams{Q: query})
	if err != nil {
		switch err.(type) {
		case *search.SearchInternalServerError:
			return search.SearchOK{}, fmt.Errorf("Internal server error occurred while searching projects")
		default:
			return search.SearchOK{}, fmt.Errorf("Unknown error occurred while searching projects: %v", err)
		}
	}
	return *response, nil
}

func LogsProject(projectName string) (*project.GetLogsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize client context for fetching logs for project %s", projectName)
	}

	response, err := client.Project.GetLogs(ctx, &project.GetLogsParams{
		ProjectName: projectName,
		Context:     ctx,
	})
	if err != nil {
		switch err.(type) {
		case *project.GetLogsBadRequest:
			return nil, fmt.Errorf("Bad request while fetching logs for project %s", projectName)
		case *project.GetLogsInternalServerError:
			return nil, fmt.Errorf("Internal server error occurred while fetching logs for project %s", projectName)
		case *project.GetLogsUnauthorized:
			return nil, fmt.Errorf("Unauthorized to fetch logs for project %s", projectName)
		default:
			return nil, fmt.Errorf("Unknown error occurred while fetching logs for project %s: %v", projectName, err)
		}
	}

	return response, nil
}
