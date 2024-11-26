package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/statistic"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/systeminfo"
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

func GetSystemInfo() (*systeminfo.GetSystemInfoOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Systeminfo.GetSystemInfo(
		ctx,
		&systeminfo.GetSystemInfoParams{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetSystemVolumes() (*systeminfo.GetVolumesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Systeminfo.GetVolumes(
		ctx,
		&systeminfo.GetVolumesParams{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
