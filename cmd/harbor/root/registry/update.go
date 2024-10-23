package registry

import (
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/registry/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateRegistryCommand() *cobra.Command {
	var opts models.Registry

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update registry",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var registryId int64

			if len(args) > 0 {
				registryId, err = strconv.ParseInt(args[0], 10, 64)
			} else {
				registryId = prompt.GetRegistryNameFromUser()
			}
			if err != nil {
				log.Errorf("failed to parse registry id: %v", err)
			}

			opts = *api.GetRegistryResponse(registryId)
			updateView := &models.Registry{
				Name:        opts.Name,
				Type:        opts.Type,
				Description: opts.Description,
				URL:         opts.URL,
				Credential: &models.RegistryCredential{
					AccessKey:    opts.Credential.AccessKey,
					Type:         opts.Credential.Type,
					AccessSecret: opts.Credential.AccessSecret,
				},
				Insecure: opts.Insecure,
			}

			update.UpdateRegistryView(updateView)
			err = api.UpdateRegistry(updateView, registryId)
			if err != nil {
				log.Errorf("failed to update registry: %v", err)
			}
		},
	}

	return cmd
}
