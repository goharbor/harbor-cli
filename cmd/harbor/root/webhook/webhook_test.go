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
	"bytes"
	"testing"

	webhookcreate "github.com/goharbor/harbor-cli/pkg/views/webhook/create"
	webhookedit "github.com/goharbor/harbor-cli/pkg/views/webhook/edit"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/webhook"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWebhookCmd_AcceptsAndNormalizesEndpointURL(t *testing.T) {
	originalCreateWebhook := createWebhook
	t.Cleanup(func() {
		createWebhook = originalCreateWebhook
	})

	createWebhook = func(opts *webhookcreate.CreateView) error {
		assert.Equal(t, "https://example.com/webhook", opts.EndpointURL)
		return nil
	}

	cmd := CreateWebhookCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{
		"my-webhook",
		"--project", "my-project",
		"--notify-type", "http",
		"--event-type", "PUSH_ARTIFACT",
		"--endpoint-url", "example.com/webhook",
	})

	err := cmd.Execute()
	require.NoError(t, err)
}

func TestEditWebhookCmd_AcceptsAndNormalizesEndpointURL(t *testing.T) {
	originalUpdateWebhook := updateWebhook
	t.Cleanup(func() {
		updateWebhook = originalUpdateWebhook
	})

	updateWebhook = func(opts *webhookedit.EditView) error {
		assert.Equal(t, "https://example.com/webhook", opts.EndpointURL)
		return nil
	}

	cmd := EditWebhookCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{
		"--project", "my-project",
		"--webhook-id", "1",
		"--notify-type", "http",
		"--event-type", "PUSH_ARTIFACT",
		"--endpoint-url", "example.com/webhook",
	})

	err := cmd.Execute()
	require.NoError(t, err)
}

func TestEditWebhookCmd_InvalidEndpointURLReturnsValidationError(t *testing.T) {
	cmd := EditWebhookCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{
		"--project", "my-project",
		"--webhook-id", "1",
		"--notify-type", "http",
		"--event-type", "PUSH_ARTIFACT",
		"--endpoint-url", "http://",
	})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "URL must contain a valid host")
}

func TestListWebhookCmd_OutputFormatSkipsTableView(t *testing.T) {
	originalListWebhooks := listWebhooks
	originalListWebhooksView := listWebhooksView
	t.Cleanup(func() {
		listWebhooks = originalListWebhooks
		listWebhooksView = originalListWebhooksView
	})

	listWebhooks = func(projectName string) (webhook.ListWebhookPoliciesOfProjectOK, error) {
		return webhook.ListWebhookPoliciesOfProjectOK{
			Payload: []*models.WebhookPolicy{
				{ID: 1, Name: "test-webhook"},
			},
		}, nil
	}

	viewCalled := false
	listWebhooksView = func(policies []*models.WebhookPolicy) {
		viewCalled = true
	}

	viper.Set("output-format", "json")
	t.Cleanup(func() {
		viper.Set("output-format", "")
	})

	cmd := ListWebhookCommand()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"my-project"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.False(t, viewCalled, "ListWebhooks view should not be called when output-format is specified")
}
