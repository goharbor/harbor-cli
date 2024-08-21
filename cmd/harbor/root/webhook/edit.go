package webhook

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/webhook/edit"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
)

func EditWebhookCmd() *cobra.Command {
	var opts edit.EditView
	var webhookId string
	var webhookIdInt int64
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "webhook edit",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			editView := &edit.EditView{
				WebhookId:               opts.WebhookId,
				ProjectName:             opts.ProjectName,
				Name:                    opts.Name,
				Description:             opts.Description,
				NotifyType:              opts.NotifyType,
				PayloadFormat:           opts.PayloadFormat,
				EventType:               opts.EventType,
				EndpointURL:             opts.EndpointURL,
				AuthHeader:              opts.AuthHeader,
				VerifyRemoteCertificate: opts.VerifyRemoteCertificate,
			}
			if opts.ProjectName != "" &&
				webhookId != "" &&
				opts.Name != "" &&
				opts.NotifyType != "" &&
				len(opts.EventType) != 0 &&
				opts.EndpointURL != "" {
				webhookIdInt, err = strconv.ParseInt(webhookId, 10, 64)
				opts.WebhookId = webhookIdInt
				err = api.UpdateWebhook(&opts)
			} else {
				err = editWebhookView(editView)
			}

			if err != nil {
				log.Errorf("failed to create webhook: %v", err)
			}

		},
	}
	flags := cmd.Flags()

	flags.StringVarP(&opts.ProjectName, "project", "", "", "Project Name")
	flags.StringVarP(&webhookId, "webhook-id", "", "", "Webhook ID")
	flags.StringVarP(&opts.Name, "name", "", "", "Webhook Name")
	flags.StringVarP(&opts.Description, "description", "", "", "Webhook Description")
	flags.StringVarP(&opts.NotifyType, "notify-type", "", "", "Notify Type (http, slack)")
	flags.StringArrayVarP(&opts.EventType, "event-type", "", []string{}, "Event Types (comma separated)")
	flags.StringVarP(&opts.EndpointURL, "endpoint-url", "", "", "Webhook Endpoint URL")
	flags.StringVarP(&opts.PayloadFormat, "payload-format", "", "", "Payload Format (Default, CloudEvents)")
	flags.StringVarP(&opts.AuthHeader, "auth-header", "", "", "Authentication Header")
	flags.BoolVarP(&opts.VerifyRemoteCertificate, "verify-remote-certificate", "", true, "Verify Remote Certificate")
	flags.BoolVarP(&opts.Enabled, "enabled", "", true, "Webhook Enabled")

	return cmd
}

func editWebhookView(view *edit.EditView) error {
	var selectedWebhook models.WebhookPolicy
	view.ProjectName = prompt.GetProjectNameFromUser()
	selectedWebhook = prompt.GetWebhookFromUser(view.ProjectName)

	view.WebhookId = selectedWebhook.ID
	view.Description = selectedWebhook.Description
	view.Enabled = selectedWebhook.Enabled
	view.EventType = selectedWebhook.EventTypes
	view.Name = selectedWebhook.Name

	view.EndpointURL = selectedWebhook.Targets[0].Address
	view.AuthHeader = selectedWebhook.Targets[0].AuthHeader
	view.PayloadFormat = string(selectedWebhook.Targets[0].PayloadFormat)
	view.VerifyRemoteCertificate = !selectedWebhook.Targets[0].SkipCertVerify
	view.NotifyType = selectedWebhook.Targets[0].Type

	edit.WebhookEditView(view)
	return api.UpdateWebhook(view)
}
