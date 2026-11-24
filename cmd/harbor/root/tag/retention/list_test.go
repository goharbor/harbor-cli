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
package retention

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListRetentionRulesCommand_ConflictingFlags(t *testing.T) {
	cmd := ListRetentionRulesCommand()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--project-id", "1", "--project-name", "my-project"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot specify both --project-id and --project-name flags")
}

func TestListRetentionRulesCommand_OutputFormat(t *testing.T) {
	policy := &models.RetentionPolicy{
		Algorithm: "or",
		ID:        42,
		Rules: []*models.RetentionRule{
			{
				Action:   "retain",
				Priority: 1,
				Template: "always",
			},
		},
	}

	tests := []struct {
		name          string
		format        string
		expectError   bool
		expectedTexts []string
	}{
		{
			name:        "JSON output formats payload directly without envelope",
			format:      "json",
			expectError: false,
			expectedTexts: []string{
				`"algorithm": "or"`,
				`"id": 42`,
				`"action": "retain"`,
			},
		},
		{
			name:        "YAML output formats payload as YAML",
			format:      "yaml",
			expectError: false,
			expectedTexts: []string{
				"algorithm: or",
				"id: 42",
				"action: retain",
			},
		},
		{
			name:        "Unsupported format returns error",
			format:      "invalid",
			expectError: true,
			expectedTexts: []string{
				"unable to output in the specified 'invalid' format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := utils.PrintFormat(policy, tt.format)
			_ = w.Close()
			_, _ = io.Copy(&buf, r)
			os.Stdout = oldStdout

			if tt.expectError {
				require.Error(t, err)
				for _, text := range tt.expectedTexts {
					assert.Contains(t, err.Error(), text)
				}
			} else {
				require.NoError(t, err)
				output := buf.String()
				// Verify raw envelope "Payload" is not present
				assert.NotContains(t, output, `"Payload":`)
				for _, text := range tt.expectedTexts {
					assert.Contains(t, output, text)
				}
			}
		})
	}
}
