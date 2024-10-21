package registry

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

				Name:        opts.Name,
				Type:        opts.Type,
				Description: opts.Description,
				URL:         opts.URL,
					AccessKey:    opts.Credential.AccessKey,
					Type:         opts.Credential.Type,
					AccessSecret: opts.Credential.AccessSecret,
				},
				Insecure: opts.Insecure,
			}

			if err != nil {
				log.Errorf("failed to update registry: %v", err)
			}
		},
	}

}
