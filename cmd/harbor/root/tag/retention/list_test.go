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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/retention"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleRetentionOK() retention.GetRetentionOK {
	return retention.GetRetentionOK{
		Payload: &models.RetentionPolicy{
			Algorithm: "or",
			ID:        42,
			Rules: []*models.RetentionRule{
				{
					Action:   "retain",
					Priority: 1,
					Template: "always",
				},
			},
		},
	}
}

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

func TestFormatRetentionRules(t *testing.T) {
	resp := sampleRetentionOK()

	tests := []struct {
		name          string
		format        string
		expectError   bool
		expectedTexts []string
	}{
		{
			name:        "JSON formats resp.Payload directly without Swagger envelope",
			format:      "json",
			expectError: false,
			expectedTexts: []string{
				`"algorithm": "or"`,
				`"id": 42`,
				`"action": "retain"`,
			},
		},
		{
			name:        "YAML formats resp.Payload as valid YAML",
			format:      "yaml",
			expectError: false,
			expectedTexts: []string{
				"algorithm: or",
				"id: 42",
				"action: retain",
			},
		},
		{
			name:        "CSV returns explicit unsupported error without panic",
			format:      "csv",
			expectError: true,
			expectedTexts: []string{
				"csv output is not supported for tag retention rules",
			},
		},
		{
			name:        "Unsupported format returns formatting error",
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

			err := formatRetentionRules(resp, tt.format)
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
				// Verify Swagger envelope {"Payload": ...} is not present
				assert.NotContains(t, output, `"Payload":`)
				for _, text := range tt.expectedTexts {
					assert.Contains(t, output, text)
				}
			}
		})
	}
}

func TestListRetentionRulesCommand_Execution(t *testing.T) {
	origGetRetentionId := getRetentionId
	origListRetention := listRetention
	t.Cleanup(func() {
		getRetentionId = origGetRetentionId
		listRetention = origListRetention
		viper.Set("output-format", "")
	})

	getRetentionId = func(projectIDStr string, isName bool) (string, error) {
		assert.Equal(t, "my-project", projectIDStr)
		assert.True(t, isName)
		return "10", nil
	}
	listRetention = func(retentionID string) (retention.GetRetentionOK, error) {
		assert.Equal(t, "10", retentionID)
		return sampleRetentionOK(), nil
	}

	t.Run("execute with json format", func(t *testing.T) {
		viper.Set("output-format", "json")

		var buf bytes.Buffer
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cmd := ListRetentionRulesCommand()
		cmd.SetArgs([]string{"--project-name", "my-project"})
		err := cmd.Execute()

		_ = w.Close()
		_, _ = io.Copy(&buf, r)
		os.Stdout = oldStdout

		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, `"algorithm": "or"`)
		assert.NotContains(t, output, `"Payload":`)
	})

	t.Run("execute with csv returns error", func(t *testing.T) {
		viper.Set("output-format", "csv")

		cmd := ListRetentionRulesCommand()
		cmd.SetArgs([]string{"--project-name", "my-project"})
		err := cmd.Execute()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "csv output is not supported for tag retention rules")
	})

	t.Run("execute with invalid format returns error", func(t *testing.T) {
		viper.Set("output-format", "invalid")

		cmd := ListRetentionRulesCommand()
		cmd.SetArgs([]string{"--project-name", "my-project"})
		err := cmd.Execute()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to output in the specified 'invalid' format")
	})
}
