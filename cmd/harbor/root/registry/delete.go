package registry

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteRegistryCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete registry",
		Example: "harbor registry delete [registryname]",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registryName, _ := api.GetRegistryIdByName(args[0])
				err = api.DeleteRegistry(registryName)
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
