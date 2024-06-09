package root

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
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
		UpdateConfigCmd(),
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
				log.Fatalf("failed to store the configuration: %v", err)
			}

			utils.PrintPayloadInJSONFormat(res.Payload)
		},
	}

	return cmd
}

func UpdateConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "update system configurations",
		Example: `harbor config update`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := utils.GetConfigurations()
			if err != nil {
				log.Fatalf("failed to get config from file: %v", err)
			}

			err = api.UpdateConfiguration(config)
			if err != nil {
				log.Fatalf("failed to update config: %v", err)
			}

			log.Infof("Configuration updated successfully.")
		},
	}

	return cmd
}
