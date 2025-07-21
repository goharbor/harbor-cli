package configurations

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Get System configuration command
func GetConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "get system configurations",
		Example: `harbor config get`,
		Run: func(cmd *cobra.Command, args []string) {
			res, err := api.GetConfigurations()
			if err != nil {
				log.Fatalf("Error getting configuration: %v", err)
			}

			if err = utils.AddConfigurationsToConfigFile(res.Payload); err != nil {
				log.Fatalf("failed to store the configuration: %v", err)
			}

			// utils.PrintPayloadInJSONFormat(res.Payload)
		},
	}

	return cmd
}
