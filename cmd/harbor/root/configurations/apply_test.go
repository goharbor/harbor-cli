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

package configurations

import (
	"strings"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/testutil"
	"github.com/spf13/cobra"
)

func TestApplyConfigCmd_Metadata(t *testing.T) {
	cmd := ApplyConfigCmd()

	if cmd == nil {
		t.Fatal("command should not be nil")
	}

	if cmd.Use != "apply" {
		t.Fatalf("expected command 'apply', got %s", cmd.Use)
	}

	if cmd.Short == "" {
		t.Fatal("Short description should not be empty")
	}
}

func TestApplyConfigCmd_RunExists(t *testing.T) {
	cmd := ApplyConfigCmd()

	if cmd.RunE == nil {
		t.Fatal("Run function should be defined")
	}
}

func TestApplyConfigCmd_IsCobraCommand(t *testing.T) {
	cmd := ApplyConfigCmd()

	if _, ok := interface{}(cmd).(*cobra.Command); !ok {
		t.Fatal("expected cobra command")
	}
}

func TestApplyConfigCmd_Errors(t *testing.T) {
	tests := []struct {
		name        string
		flags       []string
		expectError bool
		errContains string
	}{
		{
			name:        "no config file specified",
			flags:       []string{},
			expectError: true,
			errContains: "no config file specified",
		},
		{
			name:        "extra arguments",
			flags:       []string{"-f", "test.yaml", "extra-arg"},
			expectError: true,
		},
		{
			name:        "unknown flag",
			flags:       []string{"--invalid-flag"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testutil.TestCmd(t, ApplyConfigCmd, tt.flags...)

			if tt.expectError && err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.errContains != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}
