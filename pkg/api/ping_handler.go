package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/ping"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
)

func Ping() error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		logrus.Errorf("failed to get client: %v", err)
		return err
	}

	_, err = client.Ping.GetPing(ctx, &ping.GetPingParams{})
	if err != nil {
		logrus.Errorf("failed to ping the server: %v", err)
		return err
	}

	return nil
}
