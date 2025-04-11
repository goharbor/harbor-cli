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
package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/webhook"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/webhook/create"
	"github.com/goharbor/harbor-cli/pkg/views/webhook/edit"
	log "github.com/sirupsen/logrus"
)

func ListWebhooks(projectName string) (webhook.ListWebhookPoliciesOfProjectOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return webhook.ListWebhookPoliciesOfProjectOK{}, err
	}

	response, err := client.Webhook.ListWebhookPoliciesOfProject(ctx, &webhook.ListWebhookPoliciesOfProjectParams{
		ProjectNameOrID: projectName,
	})

	if err != nil {
		return webhook.ListWebhookPoliciesOfProjectOK{}, err
	}
	return *response, nil
}

func CreateWebhook(opts *create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.Webhook.CreateWebhookPolicyOfProject(ctx, &webhook.CreateWebhookPolicyOfProjectParams{
		ProjectNameOrID: opts.ProjectName,
		Policy: &models.WebhookPolicy{
			Description: opts.Description,
			Enabled:     true,
			EventTypes:  opts.EventType,
			Name:        opts.Name,
			Targets: []*models.WebhookTargetObject{
				{
					Address:        opts.EndpointURL,
					AuthHeader:     opts.AuthHeader,
					PayloadFormat:  models.PayloadFormatType(opts.PayloadFormat),
					SkipCertVerify: !opts.VerifyRemoteCertificate,
					Type:           opts.NotifyType,
				},
			},
		},
	})

	if err != nil {
		log.Errorf("%s", err)
		return err
	}

	if response != nil {
		log.Infof("Webhook `%s` created successfully", opts.Name)
	}

	return nil
}

func DeleteWebhook(projectName string, webhookId int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		log.Errorf("%s", err)
		return err
	}
	response, err := client.Webhook.DeleteWebhookPolicyOfProject(ctx, &webhook.DeleteWebhookPolicyOfProjectParams{
		WebhookPolicyID: webhookId,
		ProjectNameOrID: projectName,
	})
	if err != nil {
		log.Errorf("%s", err)
		return err
	}
	if response != nil {
		log.Infof("Webhook Id:`%d` deleted successfully", webhookId)
	}
	return nil
}

func UpdateWebhook(opts *edit.EditView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		log.Errorf("%s", err)
		return err
	}

	response, err := client.Webhook.UpdateWebhookPolicyOfProject(ctx, &webhook.UpdateWebhookPolicyOfProjectParams{
		ProjectNameOrID: opts.ProjectName,
		WebhookPolicyID: opts.WebhookId,
		Policy: &models.WebhookPolicy{
			Description: opts.Description,
			Enabled:     opts.Enabled,
			EventTypes:  opts.EventType,
			Name:        opts.Name,
			Targets: []*models.WebhookTargetObject{
				{
					Address:        opts.EndpointURL,
					AuthHeader:     opts.AuthHeader,
					PayloadFormat:  models.PayloadFormatType(opts.PayloadFormat),
					SkipCertVerify: !opts.VerifyRemoteCertificate,
					Type:           opts.NotifyType,
				},
			},
		},
	})
	if err != nil {
		log.Errorf("%s", err)
		return err
	}
	if response != nil {
		log.Infof("Webhook Id:`%d` Updated successfully", opts.WebhookId)
	}
	return nil
}
