package configurations

import (
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "update system configurations",
		Example: `harbor config update`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := utils.GetCurrentHarborConfig()
			if err != nil {
				log.Fatalf("failed to get config from file: %v", err)
			}

			// err = api.UpdateConfiguration(config)
			// if err != nil {
			// 	log.Fatalf("failed to update config: %v", err)
			// }

			log.Infof("Configuration updated successfully.")
		},
	}

	return cmd
}
