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
package policies

import (
	"testing"

	"github.com/goharbor/harbor-cli/pkg/views/replication/policies/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasReplicationUpdateFlagChanges_DestNamespaceFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "no flags returns false",
			args:     []string{},
			expected: false,
		},
		{
			name:     "dest-namespace flag triggers non-interactive mode",
			args:     []string{"--dest-namespace", "new-ns"},
			expected: true,
		},
		{
			name:     "dest-namespace-replace-count flag triggers non-interactive mode",
			args:     []string{"--dest-namespace-replace-count", "2"},
			expected: true,
		},
		{
			name:     "both dest-namespace flags trigger non-interactive mode",
			args:     []string{"--dest-namespace", "new-ns", "--dest-namespace-replace-count", "2"},
			expected: true,
		},
		{
			name:     "other flags still trigger non-interactive mode",
			args:     []string{"--name", "new-name"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := UpdateCommand()
			cmd.SetArgs(tt.args)
			// Parse the command to populate flags
			require.NoError(t, cmd.ParseFlags(tt.args))

			got := hasReplicationUpdateFlagChanges(cmd)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestApplyReplicationUpdateFlags_DestNamespace_ExplicitOverlay(t *testing.T) {
	tests := []struct {
		name                  string
		initialDestNamespace  string
		initialReplaceCount   int8
		args                  []string
		expectedDestNamespace string
		expectedReplaceCount  int8
		expectError           bool
		errorContains         string
	}{
		{
			name:                  "dest-namespace flag updates the value",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   0,
			args:                  []string{"--dest-namespace", "new-ns"},
			expectedDestNamespace: "new-ns",
			expectedReplaceCount:  0, // unchanged
			expectError:           false,
		},
		{
			name:                  "dest-namespace-replace-count flag updates the value",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   0,
			args:                  []string{"--dest-namespace-replace-count", "2"},
			expectedDestNamespace: "old-ns", // unchanged
			expectedReplaceCount:  2,
			expectError:           false,
		},
		{
			name:                  "both flags update their respective values",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   0,
			args:                  []string{"--dest-namespace", "new-ns", "--dest-namespace-replace-count", "1"},
			expectedDestNamespace: "new-ns",
			expectedReplaceCount:  1,
			expectError:           false,
		},
		{
			name:                  "dest-namespace-replace-count with value 0 is valid",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   2,
			args:                  []string{"--dest-namespace-replace-count", "0"},
			expectedDestNamespace: "old-ns",
			expectedReplaceCount:  0,
			expectError:           false,
		},
		{
			name:                  "dest-namespace-replace-count with value 1 is valid",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   2,
			args:                  []string{"--dest-namespace-replace-count", "1"},
			expectedDestNamespace: "old-ns",
			expectedReplaceCount:  1,
			expectError:           false,
		},
		{
			name:                  "dest-namespace-replace-count with value 3 is valid",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   2,
			args:                  []string{"--dest-namespace-replace-count", "3"},
			expectedDestNamespace: "old-ns",
			expectedReplaceCount:  3,
			expectError:           false,
		},
		{
			name:                  "dest-namespace-replace-count with value -1 is valid",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   2,
			args:                  []string{"--dest-namespace-replace-count", "-1"},
			expectedDestNamespace: "old-ns",
			expectedReplaceCount:  -1,
			expectError:           false,
		},
		{
			name:                  "dest-namespace-replace-count with value -2 is invalid",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   2,
			args:                  []string{"--dest-namespace-replace-count", "-2"},
			expectedDestNamespace: "old-ns",
			expectedReplaceCount:  2, // unchanged due to error
			expectError:           true,
			errorContains:         "cannot be less than -1",
		},
		{
			name:                  "dest-namespace-replace-count with value 4 is invalid",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   2,
			args:                  []string{"--dest-namespace-replace-count", "4"},
			expectedDestNamespace: "old-ns",
			expectedReplaceCount:  2, // unchanged due to error
			expectError:           true,
			errorContains:         "cannot be less than -1 or greater than 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := UpdateCommand()
			createView := &create.CreateView{
				Name:                      "test-policy",
				DestNamespace:             tt.initialDestNamespace,
				DestNamespaceReplaceCount: tt.initialReplaceCount,
				ReplicationMode:           "Push",
			}

			// Build the opts from the command args
			require.NoError(t, cmd.ParseFlags(tt.args))

			// Extract opts from the command
			opts := updateOpts{}
			flags := cmd.Flags()

			if flags.Changed("dest-namespace") {
				opts.DestNamespace, _ = flags.GetString("dest-namespace")
			}
			if flags.Changed("dest-namespace-replace-count") {
				opts.DestNamespaceReplaceCount, _ = flags.GetInt8("dest-namespace-replace-count")
			}

			err := applyReplicationUpdateFlags(cmd, createView, opts)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedDestNamespace, createView.DestNamespace)
			assert.Equal(t, tt.expectedReplaceCount, createView.DestNamespaceReplaceCount)
		})
	}
}

func TestApplyReplicationUpdateFlags_DestNamespace_PreservationWhenOmitted(t *testing.T) {
	tests := []struct {
		name                  string
		initialDestNamespace  string
		initialReplaceCount   int8
		args                  []string
		expectedDestNamespace string
		expectedReplaceCount  int8
	}{
		{
			name:                  "no dest-namespace flags preserves existing values",
			initialDestNamespace:  "existing-ns",
			initialReplaceCount:   2,
			args:                  []string{"--description", "new description"}, // other flag only
			expectedDestNamespace: "existing-ns",
			expectedReplaceCount:  2,
		},
		{
			name:                  "empty args preserves existing values",
			initialDestNamespace:  "existing-ns",
			initialReplaceCount:   1,
			args:                  []string{},
			expectedDestNamespace: "existing-ns",
			expectedReplaceCount:  1,
		},
		{
			name:                  "other flags preserve dest-namespace values",
			initialDestNamespace:  "my-ns",
			initialReplaceCount:   0,
			args:                  []string{"--description", "new description", "--enabled"},
			expectedDestNamespace: "my-ns",
			expectedReplaceCount:  0,
		},
		{
			name:                  "dest-namespace provided but not replace-count preserves replace-count",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   3,
			args:                  []string{"--dest-namespace", "new-ns"},
			expectedDestNamespace: "new-ns",
			expectedReplaceCount:  3, // preserved
		},
		{
			name:                  "dest-namespace-replace-count provided but not namespace preserves namespace",
			initialDestNamespace:  "old-ns",
			initialReplaceCount:   1,
			args:                  []string{"--dest-namespace-replace-count", "2"},
			expectedDestNamespace: "old-ns", // preserved
			expectedReplaceCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := UpdateCommand()
			createView := &create.CreateView{
				Name:                      "test-policy",
				DestNamespace:             tt.initialDestNamespace,
				DestNamespaceReplaceCount: tt.initialReplaceCount,
				ReplicationMode:           "Push",
			}

			require.NoError(t, cmd.ParseFlags(tt.args))

			// Extract opts from the command
			opts := updateOpts{}
			flags := cmd.Flags()

			if flags.Changed("dest-namespace") {
				opts.DestNamespace, _ = flags.GetString("dest-namespace")
			}
			if flags.Changed("dest-namespace-replace-count") {
				opts.DestNamespaceReplaceCount, _ = flags.GetInt8("dest-namespace-replace-count")
			}

			err := applyReplicationUpdateFlags(cmd, createView, opts)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedDestNamespace, createView.DestNamespace)
			assert.Equal(t, tt.expectedReplaceCount, createView.DestNamespaceReplaceCount)
		})
	}
}
