package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/registry/create"
	log "github.com/sirupsen/logrus"
)

func ListRegistries(opts ...ListFlags) (*registry.ListRegistriesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags

	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Registry.ListRegistries(ctx, &registry.ListRegistriesParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Q:        &listFlags.Q,
		Name:     &listFlags.Name,
		Sort:     &listFlags.Sort,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func CreateRegistry(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Registry.CreateRegistry(ctx, &registry.CreateRegistryParams{Registry: &models.Registry{Credential: &models.RegistryCredential{AccessKey: opts.Credential.AccessKey, AccessSecret: opts.Credential.AccessSecret, Type: opts.Credential.Type}, Description: opts.Description, Insecure: opts.Insecure, Name: opts.Name, Type: opts.Type, URL: opts.URL}})

	if err != nil {
		return err
	}

	log.Infof("Registry %s created", opts.Name)
	return nil
}

func DeleteRegistry(registryName int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Registry.DeleteRegistry(ctx, &registry.DeleteRegistryParams{ID: registryName})

	if err != nil {
		return err
	}

	log.Info("registry deleted successfully")

	return nil
}

func InfoRegistry(registryId int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.Registry.GetRegistry(ctx, &registry.GetRegistryParams{ID: registryId})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.Payload)
	return nil
}

func GetRegistry(registryId int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	response, err := client.Registry.GetRegistry(ctx, &registry.GetRegistryParams{ID: registryId})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.GetPayload())
	return nil
}

func UpdateRegistry(updateView *create.CreateView, projectID int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	registryUpdate := &models.RegistryUpdate{
		Name:           &updateView.Name,
		Description:    &updateView.Description,
		URL:            &updateView.URL,
		AccessKey:      &updateView.Credential.AccessKey,
		AccessSecret:   &updateView.Credential.AccessSecret,
		CredentialType: &updateView.Credential.Type,
		Insecure:       &updateView.Insecure,
	}

	_, err = client.Registry.UpdateRegistry(ctx, &registry.UpdateRegistryParams{ID: projectID, Registry: registryUpdate})

	if err != nil {
		return err
	}

	log.Info("registry updated successfully")

	return nil
}
