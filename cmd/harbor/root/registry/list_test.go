package registry

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRegistryCommand_InvalidPage(t *testing.T) {
	cmd := ListRegistryCommand()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--page", "0"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "page number must be greater than or equal to 1")
}
