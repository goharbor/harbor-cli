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

package policy

import (
	"testing"

	"github.com/goharbor/harbor-cli/pkg/testutil"
)

// TestUpdatePolicyCommand_Errors tests the custom CLI-level validations in the
// update command. API responses are not tested here.
func TestUpdatePolicyCommand_Errors(t *testing.T) {
	tests := []struct {
		name        string
		flags       []string
		expectError bool
	}{
		{
			name:        "--id flag without arguments",
			flags:       []string{"--id"},
			expectError: true,
		},
		{
			name:        "too many arguments",
			flags:       []string{"project", "policy", "extra"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testutil.TestCmd(t, UpdatePolicyCommand, tt.flags...)

			if tt.expectError && err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
