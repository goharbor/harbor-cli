package configurations

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update system configurations from local config file",
		Long: `Update Harbor system configurations using the values stored in your local config file.
This will push the configurations from ~/.config/harbor-cli/config.yaml to Harbor.`,
		Example: `harbor config update`,
		RunE: func(cmd *cobra.Command, args []string) error {
			harborConfig, err := utils.GetCurrentHarborConfig()
			if err != nil {
				return fmt.Errorf("failed to get config from file: %v", err)
			}

			if harborConfig.Configurations == (models.Configurations{}) {
				return fmt.Errorf("no configurations found in config file. Run 'harbor config get' first to populate configurations")
			}

			err = api.UpdateConfigurations(harborConfig)
			if err != nil {
				return fmt.Errorf("failed to update Harbor configurations: %v", err)
			}

			log.Infof("Harbor configurations updated successfully from local config file.")
			return nil
		},
	}

	return cmd
}
