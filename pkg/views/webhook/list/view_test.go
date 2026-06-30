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
package list

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestWebhookTargetDetailsWithoutTargets(t *testing.T) {
	webhook := &models.WebhookPolicy{}

	endpointURL, notifyType, payloadFormat := webhookTargetDetails(webhook)

	assert.Equal(t, "--", endpointURL)
	assert.Equal(t, "--", notifyType)
	assert.Equal(t, "--", payloadFormat)
}

func TestWebhookTargetDetailsWithTarget(t *testing.T) {
	webhook := &models.WebhookPolicy{
		Targets: []*models.WebhookTargetObject{
			{
				Address:       "https://example.com/hook",
				Type:          "http",
				PayloadFormat: "Default",
			},
		},
	}

	endpointURL, notifyType, payloadFormat := webhookTargetDetails(webhook)

	assert.Equal(t, "https://example.com/hook", endpointURL)
	assert.Equal(t, "http", notifyType)
	assert.Equal(t, "Default", payloadFormat)
}
