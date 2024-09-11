package api

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/health"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func GetHealth() (*health.GetHealthOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Health.GetHealth(ctx,&health.GetHealthParams{})
	if err != nil {
		return nil, fmt.Errorf("error getting health status: %w", err)
	}

	return response, nil
}
