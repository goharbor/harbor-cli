package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project_metadata"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
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

func GetProject(projectName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.Project.GetProject(ctx, &project.GetProjectParams{ProjectNameOrID: projectName})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}

func DeleteProject(projectName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Project.DeleteProject(ctx, &project.DeleteProjectParams{ProjectNameOrID: projectName})

	if err != nil {
		return err
	}

	log.Info("project deleted successfully")
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

func LogsProject(projectName string) (*project.GetLogsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Project.GetLogs(ctx, &project.GetLogsParams{
		ProjectName: projectName,
		Context:     ctx,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func AddMetadata(isID bool, projectNameOrID string, metadata map[string]string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	isName := !isID
	response, err := client.ProjectMetadata.AddProjectMetadatas(ctx, &project_metadata.AddProjectMetadatasParams{Metadata: metadata, ProjectNameOrID: projectNameOrID, XIsResourceName: &isName})
	if err != nil {
		return err
	}
	if response != nil {
		log.Info("Metadata added successfully")
	}
	return nil
}

func DeleteMetadata(isID bool, projectNameOrID string, metaName []string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	isName := !isID
	for _, meta := range metaName {
		response, err := client.ProjectMetadata.DeleteProjectMetadata(ctx, &project_metadata.DeleteProjectMetadataParams{MetaName: meta, ProjectNameOrID: projectNameOrID, XIsResourceName: &isName})
		if err != nil {
			return err
		}
		if response != nil {
			log.Info("Metadata %v deleted successfully", meta)
		}
	}

	return nil
}

func ViewMetadata(isID bool, projectNameOrID string, metaName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	isName := !isID
	response, err := client.ProjectMetadata.GetProjectMetadata(ctx, &project_metadata.GetProjectMetadataParams{MetaName: metaName, ProjectNameOrID: projectNameOrID, XIsResourceName: &isName})
	if err != nil {
		return err
	}
	if response != nil {
		log.Info("Metadata: ", response.Payload)
	}

	return nil
}

func ListMetadata(isID bool, projectNameOrID string) (project_metadata.ListProjectMetadatasOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return project_metadata.ListProjectMetadatasOK{}, err
	}

	isName := !isID
	response, err := client.ProjectMetadata.ListProjectMetadatas(ctx, &project_metadata.ListProjectMetadatasParams{ProjectNameOrID: projectNameOrID, XIsResourceName: &isName})
	if err != nil {
		return project_metadata.ListProjectMetadatasOK{}, err
	}

	return *response, nil
}

func UpdateMetadata(isID bool, projectNameOrID string, metaName string, metadata map[string]string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	isName := !isID
	response, err := client.ProjectMetadata.UpdateProjectMetadata(ctx, &project_metadata.UpdateProjectMetadataParams{MetaName: metaName, Metadata: metadata, ProjectNameOrID: projectNameOrID, XIsResourceName: &isName})
	if err != nil {
		return err
	}
	if response != nil {
		log.Info("Metadata updated successfully")
	}

	return nil
}
