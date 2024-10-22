package root

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/ping"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func pingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ping",
		Short:   "Ping the Harbor server",
		Long:    `Ping the Harbor server to check if it is alive.`,
		Example: `  harbor ping`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPing()
		},
	}

	return cmd
}

func runPing() error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		logrus.Errorf("failed to get client: %v", err)
		return err
	}

	resp, err := client.Ping.GetPing(ctx, &ping.GetPingParams{})
	if err != nil {
		logrus.Errorf("failed to ping the server: %v", err)
		return err
	}

	fmt.Printf("Ping: %s\n", resp.Payload)
	return nil
}
