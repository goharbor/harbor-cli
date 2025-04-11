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
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/webhook/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateWebhookCmd() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:   "create",
		Short: "webhook create",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
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
				log.Errorf("failed to create webhook: %v", err)
			}
		},
	}
	flags := cmd.Flags()

	flags.StringVarP(&opts.ProjectName, "project", "", "", "Project Name")
	flags.StringVarP(&opts.Name, "name", "", "", "Webhook Name")
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
	view.ProjectName = prompt.GetProjectNameFromUser()
	create.WebhookCreateView(view)
	return api.CreateWebhook(view)
}
