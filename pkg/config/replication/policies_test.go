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
package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromYAMLorJSON_DestNamespaceReplaceCount(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name                 string
		content              string
		fileType             string
		expectedReplaceCount int8
		expectError          bool
	}{
		{
			name: "omitted value produces -1 (YAML)",
			content: `
name: test-policy
dest_namespace: my-ns
`,
			fileType:             "yaml",
			expectedReplaceCount: -1,
			expectError:          false,
		},
		{
			name: "omitted value produces -1 (JSON)",
			content: `{
  "name": "test-policy",
  "dest_namespace": "my-ns"
}`,
			fileType:             "json",
			expectedReplaceCount: -1,
			expectError:          false,
		},
		{
			name: "explicit value 0 is preserved (YAML)",
			content: `
name: test-policy
dest_namespace: my-ns
dest_namespace_replace_count: 0
`,
			fileType:             "yaml",
			expectedReplaceCount: 0,
			expectError:          false,
		},
		{
			name: "explicit value 0 is preserved (JSON)",
			content: `{
  "name": "test-policy",
  "dest_namespace": "my-ns",
  "dest_namespace_replace_count": 0
}`,
			fileType:             "json",
			expectedReplaceCount: 0,
			expectError:          false,
		},
		{
			name: "explicit value 1 is preserved (YAML)",
			content: `
name: test-policy
dest_namespace: my-ns
dest_namespace_replace_count: 1
`,
			fileType:             "yaml",
			expectedReplaceCount: 1,
			expectError:          false,
		},
		{
			name: "explicit value 2 is preserved (YAML)",
			content: `
name: test-policy
dest_namespace: my-ns
dest_namespace_replace_count: 2
`,
			fileType:             "yaml",
			expectedReplaceCount: 2,
			expectError:          false,
		},
		{
			name: "explicit value 3 is preserved (YAML)",
			content: `
name: test-policy
dest_namespace: my-ns
dest_namespace_replace_count: 3
`,
			fileType:             "yaml",
			expectedReplaceCount: 3,
			expectError:          false,
		},
		{
			name: "explicit value -1 is preserved (YAML)",
			content: `
name: test-policy
dest_namespace: my-ns
dest_namespace_replace_count: -1
`,
			fileType:             "yaml",
			expectedReplaceCount: -1,
			expectError:          false,
		},
		{
			name: "explicit value -1 is preserved (JSON)",
			content: `{
  "name": "test-policy",
  "dest_namespace": "my-ns",
  "dest_namespace_replace_count": -1
}`,
			fileType:             "json",
			expectedReplaceCount: -1,
			expectError:          false,
		},
		{
			name: "value -2 is invalid (YAML)",
			content: `
name: test-policy
dest_namespace: my-ns
dest_namespace_replace_count: -2
`,
			fileType:             "yaml",
			expectedReplaceCount: 0,
			expectError:          true,
		},
		{
			name: "value 4 is invalid (YAML)",
			content: `
name: test-policy
dest_namespace: my-ns
dest_namespace_replace_count: 4
`,
			fileType:             "yaml",
			expectedReplaceCount: 0,
			expectError:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			ext := ".yaml"
			if tt.fileType == "json" {
				ext = ".json"
			}
			filename := filepath.Join(tmpDir, "test"+ext)
			err := os.WriteFile(filename, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Load config
			opts, err := LoadConfigFromFile(filename)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedReplaceCount, opts.DestNamespaceReplaceCount)
		})
	}
}

func TestLoadConfigFromYAMLorJSON_DestNamespaceReplaceCount_Distinction(t *testing.T) {
	// This test specifically verifies the distinction between:
	// 1. Omitted field (nil) -> should map to -1
	// 2. Explicit field with value 0 -> should map to 0

	tmpDir := t.TempDir()

	// Case 1: Omitted field
	omittedFile := filepath.Join(tmpDir, "omitted.yaml")
	omittedContent := `
name: test-policy-omitted
dest_namespace: my-ns
`
	err := os.WriteFile(omittedFile, []byte(omittedContent), 0644)
	require.NoError(t, err)

	optsOmitted, err := LoadConfigFromFile(omittedFile)
	require.NoError(t, err)
	assert.Equal(t, int8(-1), optsOmitted.DestNamespaceReplaceCount,
		"Omitted dest_namespace_replace_count should map to -1")

	// Case 2: Explicit field with value 0
	explicitZeroFile := filepath.Join(tmpDir, "explicit_zero.yaml")
	explicitZeroContent := `
name: test-policy-explicit-zero
dest_namespace: my-ns
dest_namespace_replace_count: 0
`
	err = os.WriteFile(explicitZeroFile, []byte(explicitZeroContent), 0644)
	require.NoError(t, err)

	optsExplicitZero, err := LoadConfigFromFile(explicitZeroFile)
	require.NoError(t, err)
	assert.Equal(t, int8(0), optsExplicitZero.DestNamespaceReplaceCount,
		"Explicit dest_namespace_replace_count: 0 should be preserved as 0")

	// Verify they are different
	assert.NotEqual(t, optsOmitted.DestNamespaceReplaceCount, optsExplicitZero.DestNamespaceReplaceCount,
		"Omitted and explicit 0 should produce different values")
}
