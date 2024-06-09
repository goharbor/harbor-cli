package root

import (
	"context"
	"fmt"
	"log"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/login"
	"github.com/spf13/cobra"
)

func ConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Manage system configurations",
		Long:    "Manage system configurations",
		Example: `harbor config get`,
	}

	cmd.AddCommand(
		GetConfigCmd(),
	)

	return cmd
}

// Get System configuration command
func GetConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "get system configurations",
		Example: `harbor config get`,
		Run: func(cmd *cobra.Command, args []string) {
			res, err := api.GetConfiguration()
			if err != nil {
				log.Fatalf("Error getting configuration: %v", err)
			}

			if err = utils.AddConfigToConfigFile(res.Payload, utils.DefaultConfigPath); err != nil {
				log.Fatalf("failed to store the credential: %v", err)
			}

			utils.PrintPayloadInJSONFormat(res.Payload)
		},
	}

	return cmd
}
