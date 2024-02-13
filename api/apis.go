package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	cnfg "github.com/goharbor/harbor-cli/internal/pkg/config"
	log "github.com/goharbor/harbor-cli/internal/pkg/logger"
	"github.com/goharbor/harbor-cli/internal/pkg/utils"
)

var (
	listArtifactObjTable []cnfg.ListArtifactRespTable
	listArtifactObjOType []cnfg.ListArtifactRespOtype
	listProjectObjTable  []cnfg.ListProjectRespTable
	listProjectObjOtype  []cnfg.ListProjectRespOtype
	listRegistryObjTable []cnfg.ListRegistryRespTable
	getProjectObjTable   []cnfg.GetProjectRespTable
	getProjectObjOtype   []cnfg.GetProjectRespOtype
	getRegistryObjTable  []cnfg.GetRegistryRespTable
	logging              log.Logger
)

func RunListArtifact(opts cnfg.ListArtifactOptions, credentialName string, OutputType string, Wide bool) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Artifact.ListArtifacts(ctx, &artifact.ListArtifactsParams{ProjectName: opts.ProjectName, RepositoryName: opts.RepositoryName})
	if err != nil {
		logging.Err("Error occured while fetching Artifacts")
		return err
	}
	// checking 4 conditions
	if int(response.XTotalCount) == 0 {
		logging.Print("No artifacts found in " + opts.RepositoryName + " repository 😢")
		return nil
	}
	if !Wide {
		logging.Normal(fmt.Sprintf("%-20s %-18s %-18s %-18s %-12s", "ARTIFACTS", "TAGS", "SIZE", "TYPE", "PUSHTIME"))
		for i := 0; i < int(response.XTotalCount); i++ {
			listArtifactObjTable = append(listArtifactObjTable, cnfg.ListArtifactRespTable{
				Digest:   utils.ShortText(response.Payload[i].Digest, 15),
				Tag:      response.Payload[i].Tags[i].Name,
				Size:     utils.ConvertSize(response.Payload[i].Size),
				Type:     response.Payload[i].Type,
				PushTime: response.Payload[i].Tags[i].PushTime.String(),
			})
		}
	}
	if OutputType == "json" || OutputType == "yaml" && !Wide {
		for i := 0; i < int(response.XTotalCount); i++ {
			listArtifactObjOType = append(listArtifactObjOType, cnfg.ListArtifactRespOtype{
				Labels:       response.Payload[i].Labels,
				Id:           response.Payload[i].ID,
				PullTime:     response.Payload[i].PullTime.String(),
				RepositoryID: response.Payload[i].RepositoryID,
				ScanOverview: response.Payload[i].ScanOverview,
			})
			listArtifactObjTable[i].Details = listArtifactObjOType
		}
	}
	if Wide {
		if OutputType == "yaml" || OutputType == "json" {
			result := utils.PrintPayloadFormat(OutputType, response)
			logging.Info(result, "")
			return nil
		} else {
			result := utils.PrintPayloadFormat("json", response)
			logging.Info(result, "")
			return nil
		}
	}
	// Printing output
	for i := 0; i < int(response.XTotalCount); i++ {
		if !Wide && OutputType == "" {
			logging.Normal(fmt.Sprintf("%-20s %-12s %-15s %-12s %-12s", listArtifactObjTable[i].Digest, listArtifactObjTable[i].Tag, listArtifactObjTable[i].Size, listArtifactObjTable[i].Type, listArtifactObjTable[i].PushTime))
		} else {
			result := utils.PrintPayloadFormat(OutputType, response)
			logging.Info(result, "")
			return nil
		}
	}
	return nil
}

func RunCreateProject(opts cnfg.CreateProjectOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.CreateProject(ctx, &project.CreateProjectParams{Project: &models.ProjectReq{ProjectName: opts.ProjectName, Public: &opts.Public, RegistryID: &opts.RegistryID, StorageLimit: &opts.StorageLimit}})

	if err != nil {
		logging.Err("Error occured while creating Project")
		return err
	}

	if response.IsCode(200) {
		logging.Normal("Project " + opts.ProjectName + " created 🎉")
	}
	return nil
}

func RunDeleteProject(opts cnfg.DeleteProjectOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.DeleteProject(ctx, &project.DeleteProjectParams{ProjectNameOrID: opts.ProjectNameOrID})

	if err != nil {
		logging.Err("Error occured while deleting Project")
		logging.Message(response.String())
		return err
	}

	if response.IsCode(200) {
		logging.Normal("Project " + opts.ProjectNameOrID + " deleted 🎉")
	}
	return nil
}

func RunListProject(opts cnfg.ListProjectOptions, credentialName string, OutputType string, Wide bool) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.ListProjects(ctx, &project.ListProjectsParams{Name: &opts.Name, Owner: &opts.Owner, Page: &opts.Page, PageSize: &opts.PageSize, Public: &opts.Public, Q: &opts.Q, Sort: &opts.Sort})

	if err != nil {
		logging.Err("Error occured while fetching Projects")
		return err
	}
	// checking 4 conditions
	if int(response.XTotalCount) == 0 {
		logging.Print("No projects found 😢")
		return nil
	}
	if !Wide {
		logging.Normal(fmt.Sprintf("%-20s %-18s %-18s %-18s", "NAME", "ACCESS", "REPOSITORY", "CREATIONTIME"))
		for i := 0; i < int(response.XTotalCount); i++ {
			listProjectObjTable = append(listProjectObjTable, cnfg.ListProjectRespTable{
				ProjectName:  response.Payload[i].Name,
				AccessLevel:  utils.AccessCheck(response.Payload[i].Metadata.Public),
				Repositories: response.Payload[i].RepoCount,
				CreationTime: response.Payload[i].CreationTime,
			})
		}
	}
	if OutputType == "json" || OutputType == "yaml" && !Wide {
		for i := 0; i < int(response.XTotalCount); i++ {
			listProjectObjOtype = append(listProjectObjOtype, cnfg.ListProjectRespOtype{
				CVEAllowlist: response.Payload[i].CVEAllowlist,
				ProjectID:    response.Payload[i].ProjectID,
				OwnerName:    response.Payload[i].OwnerName,
				Metadata:     response.Payload[i].Metadata,
			})
			listProjectObjTable[i].Details = listProjectObjOtype
		}
	}
	if Wide {
		if OutputType == "yaml" || OutputType == "json" {
			result := utils.PrintPayloadFormat(OutputType, response)
			logging.Info(result, "")
			return nil
		} else {
			result := utils.PrintPayloadFormat("json", response)
			logging.Info(result, "")
			return nil
		}
	}
	// Printing output
	for i := 0; i < int(response.XTotalCount); i++ {
		if !Wide && OutputType == "" {
			logging.Normal(fmt.Sprintf("%-20s %-12s %-15d %-12s", listProjectObjTable[i].ProjectName, listProjectObjTable[i].AccessLevel, listProjectObjTable[i].Repositories, listProjectObjTable[i].CreationTime))
		} else {
			result := utils.PrintPayloadFormat(OutputType, response)
			logging.Info(result, "")
		}
	}
	return nil
}

func RunGetProject(opts cnfg.GetProjectOptions, credentialName string, OutputType string, Wide bool) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.GetProject(ctx, &project.GetProjectParams{ProjectNameOrID: opts.ProjectNameOrID})

	if err != nil {
		logging.Err("Error occured while fetching Project")
		return err
	}
	// checking 4 conditions
	if response.IsCode(404) {
		logging.Print("Project " + opts.ProjectNameOrID + " not found 😢")
		return nil
	}

	if !Wide {
		logging.Normal(fmt.Sprintf("%-20s %-18s %-18s %-18s ", "NAME", "ACCESS", "REPOSITORY", "CREATIONTIME"))
		getProjectObjTable = append(getProjectObjTable, cnfg.GetProjectRespTable{
			ProjectName:  response.Payload.Name,
			AccessLevel:  utils.AccessCheck(response.Payload.Metadata.Public),
			Repositories: response.Payload.RepoCount,
			CreationTime: response.Payload.CreationTime,
		})
	}
	if OutputType == "json" || OutputType == "yaml" && !Wide {
		getProjectObjOtype = append(getProjectObjOtype, cnfg.GetProjectRespOtype{
			CVEAllowlist: response.Payload.CVEAllowlist,
			ProjectID:    response.Payload.ProjectID,
			OwnerName:    response.Payload.OwnerName,
			Metadata:     response.Payload.Metadata,
		})
		getProjectObjTable[0].Details = getProjectObjOtype
	}
	if Wide {
		if OutputType == "yaml" || OutputType == "json" {
			result := utils.PrintPayloadFormat(OutputType, response)
			logging.Info(result, "")
			return nil
		} else {
			result := utils.PrintPayloadFormat("json", response)
			logging.Info(result, "")
			return nil
		}
	}
	// Printing output
	if !Wide && OutputType == "" {
		logging.Normal(fmt.Sprintf("%-20s %-12s %-15d %-12s", getProjectObjTable[0].ProjectName, getProjectObjTable[0].AccessLevel, getProjectObjTable[0].Repositories, getProjectObjTable[0].CreationTime))
	} else {
		result := utils.PrintPayloadFormat(OutputType, response)
		logging.Info(result, "")
		return nil
	}
	return nil
}

// registry

func RunCreateRegistry(opts cnfg.CreateRegistrytOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Registry.CreateRegistry(ctx, &registry.CreateRegistryParams{Registry: &models.Registry{Credential: &models.RegistryCredential{AccessKey: opts.Credential.AccessKey, AccessSecret: opts.Credential.AccessSecret, Type: opts.Credential.Type}, Description: opts.Description, Insecure: opts.Insecure, Name: opts.Name, Type: opts.Type, URL: opts.Url}})

	if err != nil {
		logging.Err("Error occured while creating Registry")
		return err
	}

	if response.IsCode(200) {
		logging.Normal("Registry " + opts.Name + " created 🎉")
	}
	return nil
}

func RunDeleteRegistry(opts cnfg.DeleteRegistryOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Registry.DeleteRegistry(ctx, &registry.DeleteRegistryParams{ID: opts.Id})

	if err != nil {
		logging.Err("Error occured while deleting Registry")
		logging.Message(response.String())
		return err
	}

	if response.IsCode(200) {
		logging.Normal("Project with ID " + strconv.Itoa(int(opts.Id)) + " deleted 🎉")
	}
	return nil
}

func RunListRegistry(opts cnfg.ListRegistryOptions, credentialName string, OutputType string, Wide bool) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Registry.ListRegistries(ctx, &registry.ListRegistriesParams{Page: &opts.Page, PageSize: &opts.PageSize, Q: &opts.Q, Sort: &opts.Sort})

	if err != nil {
		logging.Err("Error occured while fetching Registries")
		return err
	}
	// checking 4 conditions
	if int(response.XTotalCount) == 0 {
		logging.Print("No registries found 😢")
		return nil
	}

	if OutputType == "json" || OutputType == "yaml" {
		result := utils.PrintPayloadFormat(OutputType, response)
		logging.Info(result, "")
		return nil
	}
	if Wide {
		result := utils.PrintPayloadFormat("json", response)
		logging.Info(result, "")
		return nil
	}
	logging.Normal(fmt.Sprintf("%-20s %-18s %-18s %-18s %-12s", "NAME", "STATUS", "TYPE", "URL", "CREATIONTIME"))
	for i := 0; i < int(response.XTotalCount); i++ {
		listRegistryObjTable = append(listRegistryObjTable, cnfg.ListRegistryRespTable{
			Name:         response.Payload[i].Name,
			Status:       response.Payload[i].Status,
			Type:         response.Payload[i].Type,
			Url:          response.Payload[i].URL,
			CreationTime: response.Payload[i].CreationTime,
		})
	}
	// Printing output
	for i := 0; i < int(response.XTotalCount); i++ {
		logging.Normal(fmt.Sprintf("%-20s %-12s %-12s %-15s %-12s", listRegistryObjTable[i].Name, listRegistryObjTable[i].Status, listRegistryObjTable[i].Type, listRegistryObjTable[i].Url, listRegistryObjTable[i].CreationTime))

	}
	return nil
}

func RunUpdateRegistry(opts cnfg.UpdateRegistrytOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	registryUpdate := &models.RegistryUpdate{}

	if opts.Credential.AccessKey != "" {
		registryUpdate.AccessKey = &opts.Credential.AccessKey
	}

	if opts.Credential.AccessSecret != "" {
		registryUpdate.AccessSecret = &opts.Credential.AccessSecret
	}

	if opts.Credential.Type != "" {
		registryUpdate.CredentialType = &opts.Credential.Type
	}

	if opts.Description != "" {
		registryUpdate.Description = &opts.Description
	}

	if opts.Name != "" {
		registryUpdate.Name = &opts.Name
	}

	if opts.Url != "" {
		registryUpdate.URL = &opts.Url
	}

	registryUpdate.Insecure = &opts.Insecure

	response, err := client.Registry.UpdateRegistry(ctx, &registry.UpdateRegistryParams{ID: opts.Id, Registry: registryUpdate})

	if err != nil {
		logging.Err("Error occured while updating Registry")
		logging.Message(response.String())
		return err
	}

	if response.IsCode(200) {
		logging.Normal("Registry with ID " + strconv.Itoa(int(opts.Id)) + " updated 🎉")
	}
	return nil
}

func RunGetRegistry(opts cnfg.GetRegistryOptions, credentialName string, OutputType string, Wide bool) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Registry.GetRegistry(ctx, &registry.GetRegistryParams{ID: opts.Id})

	if err != nil {
		logging.Err("Error occured while fetching Registry")
		return err
	}
	// checking 4 conditions
	if response.IsCode(404) {
		logging.Print("Registry with ID " + strconv.Itoa(int(opts.Id)) + " not found 😢")
		return nil
	}

	if OutputType == "json" || OutputType == "yaml" {
		result := utils.PrintPayloadFormat(OutputType, response)
		logging.Info(result, "")
		return nil
	}
	if Wide {
		result := utils.PrintPayloadFormat("json", response)
		logging.Info(result, "")
		return nil
	}
	logging.Normal(fmt.Sprintf("%-20s %-18s %-18s %-18s %-12s", "NAME", "STATUS", "TYPE", "URL", "CREATIONTIME"))
	getRegistryObjTable = append(getRegistryObjTable, cnfg.GetRegistryRespTable{
		Name:         response.Payload.Name,
		Status:       response.Payload.Status,
		Type:         response.Payload.Type,
		Url:          response.Payload.URL,
		CreationTime: response.Payload.CreationTime,
	})
	// Printing output
	logging.Normal(fmt.Sprintf("%-20s %-12s %-12s %-15s %-12s", getRegistryObjTable[0].Name, getRegistryObjTable[0].Status, getRegistryObjTable[0].Type, getRegistryObjTable[0].Url, getRegistryObjTable[0].CreationTime))

	return nil
}

// repository

func RunUpdateRepository(opts cnfg.UpdateRepositoryOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	repositoryUpdate := &models.Repository{}

	if opts.Repository.Name != "" {
		repositoryUpdate.Name = opts.Repository.Name
	}
	if opts.Repository.Description != "" {
		repositoryUpdate.Description = opts.Repository.Description
	}

	if opts.Repository.ArtifactCount >= 0 {
		repositoryUpdate.ArtifactCount = opts.Repository.ArtifactCount
	}

	if opts.Repository.PullCount >= 0 {
		repositoryUpdate.PullCount = opts.Repository.PullCount
	}

	response, err := client.Repository.UpdateRepository(ctx, &repository.UpdateRepositoryParams{ProjectName: opts.ProjectName, RepositoryName: opts.RepositoryName, Repository: repositoryUpdate})

	if err != nil {
		logging.Err("Error occured while updating Repository")
		logging.Message(response.String())
		return err
	}

	if response.IsCode(200) {
		logging.Normal("Repository " + opts.RepositoryName + " updated 🎉")
	}
	return nil
}

// login

func RunLogin(opts *cnfg.LoginOptions) error {
	clientConfig := &harbor.ClientSetConfig{
		URL:      opts.ServerAddress,
		Username: opts.Username,
		Password: opts.Password,
	}
	client := utils.GetClientByConfig(clientConfig)

	ctx := context.Background()
	_, err := client.User.GetCurrentUserInfo(ctx, &user.GetCurrentUserInfoParams{})
	if err != nil {
		return fmt.Errorf("login failed, please check your credentials: %s", err)
	}

	cred := utils.Credential{
		Name:          opts.Name,
		Username:      opts.Username,
		Password:      opts.Password,
		ServerAddress: opts.ServerAddress,
	}

	if err = utils.StoreCredential(cred, true); err != nil {
		return fmt.Errorf("failed to store the credential: %s", err)
	}
	return nil
}