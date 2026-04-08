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

	"github.com/stretchr/testify/assert"
)

func TestCreateWebhookCmd_NormalizesEndpointURLBeforeValidation(t *testing.T) {
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
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create webhook")
	assert.NotContains(t, err.Error(), "invalid URL format")
	assert.NotContains(t, err.Error(), "invalid host")
}

func TestEditWebhookCmd_NormalizesEndpointURLBeforeValidation(t *testing.T) {
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
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to edit webhook")
	assert.NotContains(t, err.Error(), "invalid URL format")
	assert.NotContains(t, err.Error(), "invalid host")
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
