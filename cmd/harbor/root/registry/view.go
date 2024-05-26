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
			if len(args) > 0 {
				registryId, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Errorf("failed to parse registry id: %v", err)
				}
				api.GetRegistry(registryId)
			} else {
				registryId := utils.GetRegistryNameFromUser()
				api.GetRegistry(registryId)
			}
		},
	}

	return cmd
}
