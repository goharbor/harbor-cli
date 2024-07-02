package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/preheat"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/instance/create"
	log "github.com/sirupsen/logrus"
)

func CreateInstance(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Preheat.CreateInstance(ctx, &preheat.CreateInstanceParams{Instance: &models.Instance{Vendor: opts.Vendor,Name: opts.Name,Description: opts.Description,Endpoint: opts.Endpoint,Enabled: opts.Enabled,AuthMode: opts.AuthMode,AuthInfo: opts.AuthInfo,Insecure: opts.Insecure}})
	if err != nil {
		return err
	}

	log.Infof("Instance %s created", opts.Name)
	return nil
}

func DeleteInstance(instanceName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Preheat.DeleteInstance(ctx, &preheat.DeleteInstanceParams{PreheatInstanceName: instanceName})

	if err != nil {
		return err
	}

	log.Info("instance deleted successfully")

	return nil
}

func ListInstance(opts ...ListFlags) (*preheat.ListInstancesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags

	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Preheat.ListInstances(ctx, &preheat.ListInstancesParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Q:        &listFlags.Q,
		Sort:     &listFlags.Sort,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}