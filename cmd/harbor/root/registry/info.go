package registry

import (
	"context"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InfoRegistryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "get registry info",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registryId, _ := strconv.ParseInt(args[0], 10, 64)
				err = runInfoRegistry(registryId)
			} else {
				registryId := utils.GetRegistryNameFromUser()
				err = runInfoRegistry(registryId)
			}
			if err != nil {
				log.Errorf("failed to get registry info: %v", err)
			}

		},
	}

	return cmd
}

func runInfoRegistry(registryId int64) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Registry.GetRegistry(ctx, &registry.GetRegistryParams{ID: registryId})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.Payload)
	return nil
}
