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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/webhook/create"
	"github.com/spf13/cobra"
)

func CreateWebhookCmd() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new webhook for a Harbor project",
		Long: `This command creates a new webhook policy for a specified Harbor project.

You can either provide all required flags (project name, notify type, endpoint, etc.) directly to create the webhook non-interactively,
or leave them out and be guided through an interactive prompt to input each field. The webhook name is required as an argument.`,
		Example: `  # Create a webhook using flags
  harbor-cli webhook create my-webhook \
    --project my-project \
    --notify-type http \
    --event-type PUSH_ARTIFACT,DELETE_ARTIFACT \
    --endpoint-url https://example.com/webhook \
    --description "Webhook for artifact events" \
    --payload-format Default \
    --auth-header "Bearer mytoken"

  # Create a webhook using the interactive prompt
  harbor-cli webhook create my-webhook`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) > 0 {
				opts.Name = args[0]
			}

			createView := &create.CreateView{
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
				err = api.CreateWebhook(&opts)
			} else {
				err = createWebhookView(createView)
			}
			if err != nil {
				return fmt.Errorf("failed to create webhook: %v", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.ProjectName, "project", "", "", "Project Name")
	flags.StringVarP(&opts.Description, "description", "", "", "Webhook Description")
	flags.StringVarP(&opts.NotifyType, "notify-type", "", "", "Notify Type (http, slack)")
	flags.StringArrayVarP(&opts.EventType, "event-type", "", []string{}, "Event Types (comma separated)")
	flags.StringVarP(&opts.EndpointURL, "endpoint-url", "", "", "Webhook Endpoint URL")
	flags.StringVarP(&opts.PayloadFormat, "payload-format", "", "", "Payload Format (Default, CloudEvents)")
	flags.StringVarP(&opts.AuthHeader, "auth-header", "", "", "Authentication Header")
	flags.BoolVarP(&opts.VerifyRemoteCertificate, "verify-remote-certificate", "", true, "Verify Remote Certificate")

	return cmd
}

func createWebhookView(view *create.CreateView) error {
	if view.ProjectName == "" {
		view.ProjectName = prompt.GetProjectNameFromUser()
	}
	err := create.WebhookCreateView(view)
	if err != nil {
		return err
	}
	return api.CreateWebhook(view)
}
