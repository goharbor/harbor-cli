package registry

import (
	"context"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/registry/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewUpdateRegistryCommand creates a new `harbor update registry` command
func UpdateRegistryCommand() *cobra.Command {

	var opts models.Registry
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update registry",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var registryId int64
			credentialName := viper.GetString("current-credential-name")
			client := utils.GetClientByCredentialName(credentialName)
			ctx := context.Background()
			if len(args) > 0 {
				registryId, err = strconv.ParseInt(args[0], 10, 64)
			} else {
				registryId = utils.GetRegistryNameFromUser()
			}
			if err != nil {
				log.Errorf("failed to parse registry id: %v", err)
			}

			response, err := client.Registry.GetRegistry(ctx,&registry.GetRegistryParams{ID: registryId})
			if err != nil {
				log.Fatal(err)
			}
			opts = *response.GetPayload()
			
			
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
			err = runUpdateRegistry(updateView, registryId)
			if err != nil {
				log.Errorf("failed to update registry: %v", err)
			}
		},
	}

	return cmd
}

func runUpdateRegistry(updateView *models.Registry, projectID int64) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	registryUpdate := &models.RegistryUpdate{
		Name:           &updateView.Name,
		Description:    &updateView.Description,
		URL:            &updateView.URL,
		AccessKey:      &updateView.Credential.AccessKey,
		AccessSecret:   &updateView.Credential.AccessSecret,
		CredentialType: &updateView.Credential.Type,
		Insecure:       &updateView.Insecure,
	}

	_, err := client.Registry.UpdateRegistry(ctx, &registry.UpdateRegistryParams{ID: projectID, Registry: registryUpdate})

	if err != nil {
		return err
	}

	log.Info("registry updated successfully")

	return nil
}
