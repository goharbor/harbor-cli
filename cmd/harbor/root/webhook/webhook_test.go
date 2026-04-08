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
