package api

import (
	"context"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project_metadata"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func AddMetadata(isID bool, projectNameOrID string, metadata map[string]string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

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
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

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
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

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

func ListMetadata(isID bool, projectNameOrID string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	isName := !isID
	response, err := client.ProjectMetadata.ListProjectMetadatas(ctx, &project_metadata.ListProjectMetadatasParams{ProjectNameOrID: projectNameOrID, XIsResourceName: &isName})
	if err != nil {
		return err
	}
	if response != nil {
		log.Info("All Metadata: ", response.Payload)
	}

	return nil
}

func UpdateMetadata(isID bool, projectNameOrID string, metaName string, metadata map[string]string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

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
