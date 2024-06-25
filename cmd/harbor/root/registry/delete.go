package registry

import (
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
				err = api.DeleteRegistry(registryId)
			} else {
				registryId := prompt.GetRegistryNameFromUser()
				err = api.DeleteRegistry(registryId)
			}
			if err != nil {
				log.Errorf("failed to delete registry: %v", err)
			}
		},
	}

	return cmd
}
