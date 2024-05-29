package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/statistic"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func GetStats() (*statistic.GetStatisticOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Statistic.GetStatistic(
		ctx,
		&statistic.GetStatisticParams{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
