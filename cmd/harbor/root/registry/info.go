package registry

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func InfoRegistryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "get registry info",
		Example: "harbor registry info [registryname]",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registryName, _ := api.GetRegistryIdByName(args[0])
				err = api.InfoRegistry(registryName)
			} else {
				registryId := prompt.GetRegistryNameFromUser()
				err = api.InfoRegistry(registryId)
			}
			if err != nil {
				log.Errorf("failed to get registry info: %v", err)
			}

		},
	}

	return cmd
}
