package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/configure"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

// GetConfigurationResponse of the system
func GetConfiguration() (*configure.GetConfigurationsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Configure.GetConfigurations(
		ctx,
		&configure.GetConfigurationsParams{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Update Configuration of the system
func UpdateConfiguration(config *models.Configurations) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Configure.UpdateConfigurations(
		ctx,
		&configure.UpdateConfigurationsParams{Configurations: config},
	)
	if err != nil {
		return err
	}

	return nil
}
