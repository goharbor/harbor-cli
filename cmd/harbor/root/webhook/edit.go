package webhook

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/webhook/edit"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func EditWebhookCmd() *cobra.Command {
	var opts edit.EditView
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
				opts.Name != "" &&
				opts.NotifyType != "" &&
				len(opts.EventType) != 0 &&
				opts.EndpointURL != "" {
				err = api.UpdateWebhook(&opts)
			} else {
				err = editWebhookView(editView)
			}

			if err != nil {
				log.Errorf("failed to create webhook: %v", err)
			}

		},
	}
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
