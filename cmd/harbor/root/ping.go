package root

import (
    "fmt"
    "github.com/goharbor/harbor-cli/pkg/utils"
    "github.com/spf13/cobra"
    "github.com/sirupsen/logrus"
    "github.com/goharbor/go-client/pkg/sdk/v2.0/client/ping"
)

func GetPing() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "ping",
        Short: "Ping the Harbor API server",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, client, err := utils.ContextWithClient()
            if err != nil {
                logrus.Errorf("failed to get client: %v", err)
                return err
            }

            // Initialize params even if not used
            params := &ping.GetPingParams{}

            response, err := client.Ping.GetPing(ctx, params)
            if err != nil {
                logrus.Errorf("ping request failed: %v", err)
                return err
            }

            fmt.Printf("Ping successful: %s\n", response.Payload)
            return nil
        },
    }

    return cmd
}
