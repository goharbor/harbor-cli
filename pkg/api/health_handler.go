package api

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/health"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func GetHealth() (*health.GetHealthOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client context: ")
	}

	response, err := client.Health.GetHealth(ctx, &health.GetHealthParams{})
	if err != nil {
		switch err.(type) {
		case *health.GetHealthInternalServerError:
			return nil, fmt.Errorf("internal server error occurred while getting health status")
		default:
			return nil, fmt.Errorf("unknown error occurred while getting health status: %w", err)
		}
	}

	return response, nil
}
