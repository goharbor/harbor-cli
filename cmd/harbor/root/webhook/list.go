package webhook

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/webhook"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	webhookViews "github.com/goharbor/harbor-cli/pkg/views/webhook/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListWebhookCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list webhook",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var resp webhook.ListWebhookPoliciesOfProjectOK
			if len(args) > 0 {
				projectName := args[0]
				resp, err = api.ListWebhooks(projectName)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				resp, err = api.ListWebhooks(projectName)
			}

			if err != nil {
				log.Errorf("failed to list webhooks: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(resp)
				return
			}

			webhookViews.ListWebhooks(resp.Payload)

		},
	}
	return cmd
}
