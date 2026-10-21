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

package declarative

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportConversionsOmitSecretsAndObservedState(t *testing.T) {
	registryModel := &models.Registry{
		ID:   42,
		Name: "hub",
		Type: "docker-hub",
		URL:  "https://hub.docker.com",
		Credential: &models.RegistryCredential{
			AccessKey:    "user",
			AccessSecret: "secret",
			Type:         "basic",
		},
	}
	webhookModel := &models.WebhookPolicy{
		ID:         7,
		Name:       "notify",
		Enabled:    true,
		EventTypes: []string{"PUSH_ARTIFACT"},
		Targets: []*models.WebhookTargetObject{
			{Address: "https://example.com", AuthHeader: "secret", SkipCertVerify: true, Type: "http"},
		},
	}

	exportedRegistry := registryFromModel(registryModel)
	exportedWebhook := webhookFromModel(webhookModel)

	assert.Nil(t, exportedRegistry.Credential)
	require.Len(t, exportedWebhook.Targets, 1)
	assert.Nil(t, exportedWebhook.Targets[0].AuthHeaderFrom)
	require.NotNil(t, exportedWebhook.Targets[0].VerifyRemoteCertificate)
	assert.False(t, *exportedWebhook.Targets[0].VerifyRemoteCertificate)
}

func TestSystemConfigurationUsesAPINamesAndValues(t *testing.T) {
	response := &models.ConfigurationsResponse{
		ReadOnly:                   &models.BoolConfigItem{Value: true},
		ProjectCreationRestriction: &models.StringConfigItem{Value: "adminonly"},
	}

	configuration := systemConfiguration(response)

	assert.Equal(t, true, configuration["read_only"])
	assert.Equal(t, "adminonly", configuration["project_creation_restriction"])
	assert.NotContains(t, configuration, "ReadOnly")
}

func TestWebhookSecretReferenceIsResolvedOnlyForApply(t *testing.T) {
	t.Setenv("HARBOR_WEBHOOK_AUTH", "Bearer token")
	desired := Webhook{
		Name: "notify",
		Targets: []WebhookTarget{{
			NotifyType:     ptr("http"),
			Endpoint:       ptr("https://example.com"),
			AuthHeaderFrom: &SecretRef{Env: "HARBOR_WEBHOOK_AUTH"},
		}},
	}

	policy, err := webhookToModel(desired, nil)
	require.NoError(t, err)
	require.Len(t, policy.Targets, 1)
	assert.Equal(t, "Bearer token", policy.Targets[0].AuthHeader)
}

func TestWebhookConversionPreservesAllTargets(t *testing.T) {
	model := &models.WebhookPolicy{Targets: []*models.WebhookTargetObject{
		{Type: "http", Address: "https://one.example.com"},
		{Type: "http", Address: "https://two.example.com"},
	}}

	exported := webhookFromModel(model)
	require.Len(t, exported.Targets, 2)

	roundTrip, err := webhookToModel(exported, nil)
	require.NoError(t, err)
	require.Len(t, roundTrip.Targets, 2)
	assert.Equal(t, "https://one.example.com", roundTrip.Targets[0].Address)
	assert.Equal(t, "https://two.example.com", roundTrip.Targets[1].Address)
}

func TestWebhookCreateRequiresTarget(t *testing.T) {
	_, err := webhookToModel(Webhook{Name: "notify"}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires at least one target")
}

func TestProjectRequestWritesPublicMetadata(t *testing.T) {
	backend := &LiveBackend{projects: map[string]*models.Project{
		"demo": {Metadata: &models.ProjectMetadata{AutoScan: ptr("true")}},
	}}

	request, err := backend.projectRequest(Project{Name: "demo", Public: ptr(false)})
	require.NoError(t, err)
	require.NotNil(t, request.Metadata)
	assert.Equal(t, "false", request.Metadata.Public)
	require.NotNil(t, request.Metadata.AutoScan)
	assert.Equal(t, "true", *request.Metadata.AutoScan)
}

func TestQuotaProjectIDSupportsDecodedAPIReference(t *testing.T) {
	id, ok := quotaProjectID(map[string]any{"id": float64(23), "name": "demo"})
	assert.True(t, ok)
	assert.Equal(t, int64(23), id)
}
