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
	"testing"

	"github.com/goharbor/harbor-cli/pkg/testutil"
)

// TestCreateRetentionCommand_Errors tests the custom validations we perform in the command,
// testing of query builder and API response will be delegated to their respective places
func TestCreateRetentionCommand_Errors(t *testing.T) {
	tests := []struct {
		name        string
		flags       []string
		expectError bool
	}{
		{
			name:        "conflicting project-id and project-name flags",
			flags:       []string{"--project-id", "1", "--project-name", "test"},
			expectError: true,
		},
		{
			name:        "invalid level value",
			flags:       []string{"--level", "invalid", "--project-name", "test"},
			expectError: true,
		},
		{
			name:        "invalid action value",
			flags:       []string{"--action", "invalid", "--project-name", "test"},
			expectError: true,
		},
		{
			name:        "invalid algorithm value",
			flags:       []string{"--algorithm", "invalid", "--project-name", "test"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testutil.TestCmd(t, CreateRetentionCommand, tt.flags...)

			if tt.expectError && err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestCreateRetentionCommand_Flags tests that flags are properly defined
func TestCreateRetentionCommand_Flags(t *testing.T) {
	cmd := CreateRetentionCommand()

	// Test that required flags exist
	requiredFlags := []string{"project-name", "project-id", "level", "action", "algorithm", "taglist", "repodecoration", "repolist", "tagdecoration"}
	for _, flag := range requiredFlags {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("flag %s not found", flag)
		}
	}
}
