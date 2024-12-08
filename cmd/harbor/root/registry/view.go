package registry

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/registry/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ViewRegistryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "get registry information",
		Example: "harbor registry view [registryName]",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var registryId int64
			var registry *registry.GetRegistryOK

			if len(args) > 0 {
				registryId, err = api.GetRegistryIdByName(args[0])
				if err != nil {
					log.Errorf("failed to get registry name by id: %v", err)
					return
				}
			} else {
				registryId = prompt.GetRegistryNameFromUser()
			}

			registry, err = api.ViewRegistry(registryId)

			if err != nil {
				log.Errorf("failed to get registry info: %v", err)
				return
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(registry, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				view.ViewRegistry(registry.Payload)
			}

		},
	}

	return cmd
}
