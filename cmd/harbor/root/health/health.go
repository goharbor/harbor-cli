package health

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/health"
	"github.com/goharbor/harbor-cli/pkg/utils"
	healthview "github.com/goharbor/harbor-cli/pkg/views/health"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewListProjectCommand creates a new `harbor list project` command
func HealthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Display server health",
		Run: func(cmd *cobra.Command, args []string) {
			health, err := RunHealth()
			if err != nil {
				log.Fatalf("failed to get projects list: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag == "json" {
				utils.PrintPayloadInJSONFormat(health)
				return
			} else if FormatFlag == "yaml" {
				utils.PrintPayloadInYAMLFormat(health)
			} else {
				healthview.View(health)
			}
		},
	}

	return cmd
}

func RunHealth() (*health.GetHealthOK, error) {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Health.GetHealth(ctx, health.NewGetHealthParamsWithContext(ctx))
	if err != nil {
		return nil, err
	}
	return response, nil
}
