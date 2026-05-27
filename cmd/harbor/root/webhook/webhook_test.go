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
