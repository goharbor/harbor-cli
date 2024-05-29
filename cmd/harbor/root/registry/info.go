package registry

import (
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
				err = api.InfoRegistry(registryId)
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
