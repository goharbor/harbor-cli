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

// --- Command metadata tests ---

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

// --- YAML parsing regression tests (fixes #1023) ---
//
// These tests directly exercise parseYAMLConfig to verify that:
//   - the documented wrapped format (configurations: ...) is parsed correctly
//   - a malformed wrapped file (wrong type under "configurations") is rejected
//     with an error and does NOT silently fall back to a no-op flat parse
//   - legacy flat files continue to work for backward compatibility
//   - an empty "configurations" key is an explicit error

func TestParseYAMLConfig(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
		// checkResult is called only when wantErr is false.
		checkResult func(t *testing.T, got interface{})
	}{
		{
			name: "valid wrapped YAML",
			input: `
configurations:
  auth_mode: db_auth
`,
			wantErr: false,
			checkResult: func(t *testing.T, got interface{}) {
				if got == nil {
					t.Fatal("expected non-nil configurations")
				}
			},
		},
		{
			// Regression test for #1023: a file that starts with "configurations:"
			// but has an invalid value type must NOT silently succeed and produce
			// a no-op apply.
			name: "malformed wrapped YAML — configurations is a scalar not a map",
			input: `
configurations: "this should be a map"
`,
			wantErr:     true,
			errContains: "configurations",
		},
		{
			// Regression test for #1023: a list under configurations must be rejected.
			name: "malformed wrapped YAML — configurations is a list not a map",
			input: `
configurations:
  - invalid_list_item
`,
			wantErr:     true,
			errContains: "configurations",
		},
		{
			// Legacy flat files must continue to work (backward compatibility).
			name: "legacy flat YAML",
			input: `
auth_mode: db_auth
`,
			wantErr: false,
			checkResult: func(t *testing.T, got interface{}) {
				if got == nil {
					t.Fatal("expected non-nil configurations")
				}
			},
		},
		{
			// "configurations:" with a null/empty value must return an error,
			// not an empty configurations object that silently no-ops.
			name:        "empty configurations key in YAML",
			input:       "configurations:\n",
			wantErr:     true,
			errContains: "configurations",
		},
		{
			name:        "invalid YAML syntax",
			input:       ":\t: invalid yaml :::",
			wantErr:     true,
			errContains: "failed to parse YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseYAMLConfig([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error containing %q, got: %v", tt.errContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkResult != nil {
				tt.checkResult(t, got)
			}
		})
	}
}

// --- JSON parsing regression tests (fixes #1023) ---

func TestParseJSONConfig(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
		checkResult func(t *testing.T, got interface{})
	}{
		{
			name:    "valid wrapped JSON",
			input:   `{"configurations": {"auth_mode": "db_auth"}}`,
			wantErr: false,
			checkResult: func(t *testing.T, got interface{}) {
				if got == nil {
					t.Fatal("expected non-nil configurations")
				}
			},
		},
		{
			// Regression test for #1023: "configurations" pointing to a scalar
			// must be rejected, not silently treated as empty.
			name:        "malformed wrapped JSON — configurations is a string not an object",
			input:       `{"configurations": "not-an-object"}`,
			wantErr:     true,
			errContains: "configurations",
		},
		{
			// Regression test for #1023: "configurations" pointing to a list
			// must be rejected.
			name:        "malformed wrapped JSON — configurations is an array",
			input:       `{"configurations": ["a", "b"]}`,
			wantErr:     true,
			errContains: "configurations",
		},
		{
			// Legacy flat files must continue to work (backward compatibility).
			name:    "legacy flat JSON",
			input:   `{"auth_mode": "db_auth"}`,
			wantErr: false,
			checkResult: func(t *testing.T, got interface{}) {
				if got == nil {
					t.Fatal("expected non-nil configurations")
				}
			},
		},
		{
			// "configurations": null must return an error.
			name:        "null configurations value in JSON",
			input:       `{"configurations": null}`,
			wantErr:     true,
			errContains: "configurations",
		},
		{
			name:        "invalid JSON syntax",
			input:       `{not valid json`,
			wantErr:     true,
			errContains: "failed to parse JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJSONConfig([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error containing %q, got: %v", tt.errContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkResult != nil {
				tt.checkResult(t, got)
			}
		})
	}
}
