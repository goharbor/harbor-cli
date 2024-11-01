package registry

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// NewGetRegistryCommand creates a new `harbor get registry` command
func ViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "get registry by id",
		Example: "harbor registry view [registryname]",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registryName, err := api.GetRegistryIdByName(args[0])
				if err != nil {
					log.Errorf("failed to parse registry id: %v", err)
				}
				err = api.GetRegistry(registryName)
				if err != nil {
					log.Errorf("failed to get registry: %v", err)
				}
			} else {
				registryId := prompt.GetRegistryNameFromUser()
				err = api.GetRegistry(registryId)
			}
			if err != nil {
				log.Errorf("failed to get registry: %v", err)
			}
		},
	}

	return cmd
}
