package api

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func ListRegistries(opts ...ListFlags) (*registry.ListRegistriesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client context")
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
		switch err.(type) {
		case *registry.ListRegistriesForbidden:
			return nil, fmt.Errorf("Forbidden request while listing registries: %s", listFlags.Q)
		case *registry.ListRegistriesInternalServerError:
			return nil, fmt.Errorf("internal server error occurred while listing registries")
		case *registry.ListRegistriesUnauthorized:
			return nil, fmt.Errorf("unauthorized access to list registries")
		default:
			return nil, fmt.Errorf("unknown error occurred while listing registries: %v", err)
		}
	}

	return response, nil
}

func CreateRegistry(opts CreateRegView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client context")
	}

	_, err = client.Registry.CreateRegistry(
		ctx,
		&registry.CreateRegistryParams{
			Registry: &models.Registry{
				Credential: &models.RegistryCredential{
					AccessKey:    opts.Credential.AccessKey,
					AccessSecret: opts.Credential.AccessSecret,
					Type:         opts.Credential.Type,
				},
				Description: opts.Description,
				Insecure:    opts.Insecure,
				Name:        opts.Name,
				Type:        opts.Type,
				URL:         opts.URL,
			},
		},
	)
	if err != nil {

		switch err.(type) {
		case *registry.CreateRegistryForbidden:
			return fmt.Errorf("Forbidden request while creating registry: %s", opts.Name)
		case *registry.CreateRegistryInternalServerError:
			return fmt.Errorf("internal server error occurred while creating registry")
		case *registry.CreateRegistryUnauthorized:
			return fmt.Errorf("unauthorized access to create registry")
		case *registry.CreateRegistryBadRequest:
			return fmt.Errorf("bad request while creating registry: %s", opts.Name)
		case *registry.CreateRegistryConflict:
			return fmt.Errorf("conflict error: the registry already exists")
		default:
			return fmt.Errorf("unknown error occurred while creating registry: %v", err)
		}
	}

	log.Infof("Registry %s created", opts.Name)
	return nil
}

func DeleteRegistry(registryName int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client context")
	}
	_, err = client.Registry.DeleteRegistry(ctx, &registry.DeleteRegistryParams{ID: registryName})
	if err != nil {
		switch err.(type) {
		case *registry.DeleteRegistryNotFound:
			return fmt.Errorf("registry not found: %d", registryName)
		case *registry.DeleteRegistryForbidden:
			return fmt.Errorf("bad request while deleting registry: %d", registryName)
		case *registry.DeleteRegistryInternalServerError:
			return fmt.Errorf("internal server error occurred while deleting registry: %d", registryName)
		case *registry.DeleteRegistryUnauthorized:
			return fmt.Errorf("unauthorized access to delete registry: %d", registryName)
		default:
			return fmt.Errorf("unknown error occurred while deleting registry: %v", err)
		}
	}

	log.Info("registry deleted successfully")

	return nil
}

func ViewRegistry(registryId int64) (*registry.GetRegistryOK, error) {
	ctx, client, err := utils.ContextWithClient()
	var response = &registry.GetRegistryOK{}
	if err != nil {
		return response, fmt.Errorf("failed to initialize client context")
	}

	response, err = client.Registry.GetRegistry(ctx, &registry.GetRegistryParams{ID: registryId})

	if err != nil {
		switch err.(type) {
		case *registry.GetRegistryNotFound:
			return response, fmt.Errorf("registry not found: %d", registryId)
		case *registry.GetRegistryForbidden:
			return response, fmt.Errorf("bad request while viewing registry: %d", registryId)
		case *registry.GetRegistryInternalServerError:
			return response, fmt.Errorf("internal server error occurred while viewing registry: %d", registryId)
		case *registry.GetRegistryUnauthorized:
			return response, fmt.Errorf("unauthorized access to view registry: %d", registryId)
		default:
			return response, fmt.Errorf("unknown error occurred while viewing registry: %v", err)
		}
	}
	if response.Payload.ID == 0 {
		return response, fmt.Errorf("registry is not found")
	}

	return response, nil
}

func GetRegistryResponse(registryId int64) *models.Registry {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil
	}
	response, err := client.Registry.GetRegistry(ctx, &registry.GetRegistryParams{ID: registryId})
	if err != nil {
		return nil
	}
	if response.Payload.ID == 0 {
		return nil
	}

	return response.GetPayload()
}

func UpdateRegistry(updateView *models.Registry, projectID int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client context")
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

	_, err = client.Registry.UpdateRegistry(
		ctx,
		&registry.UpdateRegistryParams{ID: projectID, Registry: registryUpdate},
	)
	if err != nil {
		switch err.(type) {
		case *registry.UpdateRegistryNotFound:
			return fmt.Errorf("registry not found: %d", projectID)
		case *registry.UpdateRegistryForbidden:
			return fmt.Errorf("bad request while updating registry: %d", projectID)
		case *registry.UpdateRegistryInternalServerError:
			return fmt.Errorf("internal server error occurred while updating registry: %d", projectID)
		case *registry.UpdateRegistryUnauthorized:
			return fmt.Errorf("unauthorized access to update registry: %d", projectID)
		default:
			return fmt.Errorf("unknown error occurred while updating registry: %v", err)
		}
	}

	log.Info("registry updated successfully")

	return nil
}

// Get List of Registry Providers
func GetRegistryProviders() ([]string, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client context")
	}
	response, err := client.Registry.ListRegistryProviderTypes(
		ctx,
		&registry.ListRegistryProviderTypesParams{},
	)
	if err != nil {

		switch err.(type) {
		case *registry.ListRegistryProviderTypesForbidden:
			return nil, fmt.Errorf("Forbidden request while listing registry providers")
		case *registry.ListRegistryProviderTypesInternalServerError:
			return nil, fmt.Errorf("internal server error occurred while listing registry providers")
		case *registry.ListRegistryProviderTypesUnauthorized:
			return nil, fmt.Errorf("unauthorized access to list registry providers")
		default:
			return nil, fmt.Errorf("unknown error occurred while listing registry providers: %v", err)
		}

	}

	return response.Payload, nil
}

func GetRegistryIdByName(registryName string) (int64, error) {
	var opts ListFlags

	r, err := ListRegistries(opts)
	if err != nil {
		return 0, err
	}

	for _, registry := range r.Payload {
		if registry.Name == registryName {
			return registry.ID, nil
		}
	}

	return 0, err
}
