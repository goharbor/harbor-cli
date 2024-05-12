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

// NewDeleteRegistryCommand creates a new `harbor delete registry` command
func DeleteRegistryCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete registry by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registryId, _ := strconv.ParseInt(args[0], 10, 64)
				err = runDeleteRegistry(registryId)
			} else {
				registryId := utils.GetRegistryNameFromUser()
				err = runDeleteRegistry(registryId)
			}
			if err != nil {
				log.Errorf("failed to delete registry: %v", err)
			}
		},
	}

	return cmd
}

func runDeleteRegistry(registryName int64) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	_, err := client.Registry.DeleteRegistry(ctx, &registry.DeleteRegistryParams{ID: registryName})

	if err != nil {
		return err
	}

	log.Info("registry deleted successfully")

	return nil
}
