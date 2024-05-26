package registry

import (
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// NewGetRegistryCommand creates a new `harbor get registry` command
func ViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "get registry by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				registryId, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Errorf("failed to parse registry id: %v", err)
				}
				err = api.GetRegistry(registryId)
			} else {
				registryId := utils.GetRegistryNameFromUser()
				err = api.GetRegistry(registryId)
			}

			if err != nil {
				log.Errorf("failed to get registry: %v", err)
			}
		},
	}

	return cmd
}
