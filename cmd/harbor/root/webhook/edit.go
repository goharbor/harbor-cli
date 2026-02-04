// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package webhook

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/webhook/edit"
	"github.com/spf13/cobra"
)

func EditWebhookCmd() *cobra.Command {
	var opts edit.EditView
	cmd := &cobra.Command{
		Use:   "edit [WEBHOOK_NAME]",
		Short: "Edit an existing webhook for a Harbor project",
		Long: `This command allows you to update an existing webhook policy in a Harbor project.

You can either pass all the necessary flags (webhook ID, project name, etc.) to perform a non-interactive update,
or leave them out and use the interactive prompt to select and update a webhook.`,
		Example: `  # Edit a webhook by providing all fields directly
  harbor-cli webhook edit my-webhook \
    --project my-project \
    --notify-type http \
    --event-type PUSH_ARTIFACT \
    --endpoint-url https://new-url.com \
    --description "Updated webhook for artifact push" \
    --payload-format Default \
    --auth-header "Bearer newtoken" \
    --enabled=true

  # Edit a webhook using the interactive prompt
  harbor-cli webhook edit`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) > 0 {
				opts.Name = args[0]
			}

			if opts.WebhookId != -1 && opts.Name != "" {
				return fmt.Errorf("webhook ID and name cannot be provided together")
			}
			if opts.ProjectName == "" && (opts.Name != "" || opts.WebhookId != -1) {
				return fmt.Errorf("project name is required when webhook name is provided")
			}
			if opts.ProjectName != "" && opts.Name != "" {
				opts.WebhookId, err = api.GetWebhookID(opts.ProjectName, opts.Name)
				if err != nil {
					return fmt.Errorf("failed to get webhook ID: %v", err)
				}
			}

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
				opts.WebhookId != -1 &&
				opts.NotifyType != "" &&
				len(opts.EventType) != 0 &&
				opts.EndpointURL != "" {
				if err := utils.ValidateURL(opts.EndpointURL); err != nil {
					return err
				}
				err = api.UpdateWebhook(&opts)
			} else {
				err = editWebhookView(editView)
			}
			if err != nil {
				return fmt.Errorf("failed to edit webhook: %v", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.ProjectName, "project", "", "", "Project Name")
	flags.Int64VarP(&opts.WebhookId, "webhook-id", "", -1, "Webhook ID")
	flags.StringVarP(&opts.Description, "description", "", "", "Webhook Description")
	flags.StringVarP(&opts.NotifyType, "notify-type", "", "", "Notify Type (http, slack)")
	flags.StringSliceVarP(&opts.EventType, "event-type", "", []string{}, "Event Types (comma separated)")
	flags.StringVarP(&opts.EndpointURL, "endpoint-url", "", "", "Webhook Endpoint URL")
	flags.StringVarP(&opts.PayloadFormat, "payload-format", "", "", "Payload Format (Default, CloudEvents)")
	flags.StringVarP(&opts.AuthHeader, "auth-header", "", "", "Authentication Header")
	flags.BoolVarP(&opts.VerifyRemoteCertificate, "verify-remote-certificate", "", true, "Verify Remote Certificate")
	flags.BoolVarP(&opts.Enabled, "enabled", "", true, "Webhook Enabled")

	return cmd
}

func editWebhookView(view *edit.EditView) error {
	var selectedWebhook models.WebhookPolicy
	var err error
	if view.ProjectName == "" {
		view.ProjectName, err = prompt.GetProjectNameFromUser()
		if err != nil {
			return err
		}
	}
	if view.WebhookId == -1 {
		selectedWebhook, err = prompt.GetWebhookFromUser(view.ProjectName)
		if err != nil {
			return err
		}
	} else {
		selectedWebhook, err = api.GetWebhook(view.ProjectName, view.WebhookId)
		if err != nil {
			return err
		}
	}
	view.WebhookId = selectedWebhook.ID
	view.Description = selectedWebhook.Description
	view.Enabled = selectedWebhook.Enabled
	view.EventType = selectedWebhook.EventTypes
	view.Name = selectedWebhook.Name
	if len(selectedWebhook.Targets) > 0 {
		view.EndpointURL = selectedWebhook.Targets[0].Address
		view.AuthHeader = selectedWebhook.Targets[0].AuthHeader
		view.PayloadFormat = string(selectedWebhook.Targets[0].PayloadFormat)
		view.VerifyRemoteCertificate = !selectedWebhook.Targets[0].SkipCertVerify
		view.NotifyType = selectedWebhook.Targets[0].Type
	}
	edit.WebhookEditView(view)
	return api.UpdateWebhook(view)
}
