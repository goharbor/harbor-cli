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
package labels

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateLabelCommand_Structure(t *testing.T) {
	cmd := CreateLabelCommand()

	assert.Equal(t, "create", cmd.Use)
	assert.Equal(t, "create label", cmd.Short)
	assert.Equal(t, "create label in harbor", cmd.Long)
	assert.Equal(t, "harbor label create", cmd.Example)
}

func TestCreateLabelCommand_Flags(t *testing.T) {
	cmd := CreateLabelCommand()
	flags := cmd.Flags()

	tests := []struct{ name, shorthand, defValue string }{
		{"name", "n", ""},
		{"color", "", "#FFFFFF"},
		{"scope", "s", "g"},
		{"project", "i", "0"},
		{"description", "d", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := flags.Lookup(tt.name)
			assert.NotNil(t, flag, "Flag %s should exist", tt.name)
			assert.Equal(t, tt.shorthand, flag.Shorthand)
			assert.Equal(t, tt.defValue, flag.DefValue)
		})
	}
}

func TestCreateLabelCommand_ProjectScope_MissingProjectID(t *testing.T) {
	cmd := CreateLabelCommand()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{})

	assert.NoError(t, cmd.Flags().Set("name", "test-label"))
	assert.NoError(t, cmd.Flags().Set("scope", "p"))
	// project ID is 0 (default) - should trigger error before API call

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project ID is required when scope is 'p'")
	assert.Contains(t, err.Error(), "Use --project flag to specify the project ID")
}

func TestCreateLabelCommand_ExactArgs(t *testing.T) {
	cmd := CreateLabelCommand()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"extra-arg"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 0 arg(s)")
}

func TestCreateLabelCommand_DefaultValues(t *testing.T) {
	cmd := CreateLabelCommand()
	flags := cmd.Flags()

	color, err := flags.GetString("color")
	assert.NoError(t, err)
	assert.Equal(t, "#FFFFFF", color)

	scope, err := flags.GetString("scope")
	assert.NoError(t, err)
	assert.Equal(t, "g", scope)

	projectID, err := flags.GetInt64("project")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), projectID)

	name, err := flags.GetString("name")
	assert.NoError(t, err)
	assert.Equal(t, "", name)

	description, err := flags.GetString("description")
	assert.NoError(t, err)
	assert.Equal(t, "", description)
}

func TestCreateLabelCommand_ColorFlag(t *testing.T) {
	cmd := CreateLabelCommand()
	flags := cmd.Flags()

	testColors := []string{"#000000", "#FFFFFF", "#FF5501", "#48960C"}
	for _, color := range testColors {
		t.Run(color, func(t *testing.T) {
			assert.NoError(t, flags.Set("color", color))
			val, err := flags.GetString("color")
			assert.NoError(t, err)
			assert.Equal(t, color, val)
		})
	}
}

func TestCreateLabelCommand_ScopeValues(t *testing.T) {
	tests := []struct{ scope, expected string }{
		{"g", "g"},
		{"p", "p"},
	}
	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			cmd := CreateLabelCommand()
			assert.NoError(t, cmd.Flags().Set("scope", tt.scope))
			val, err := cmd.Flags().GetString("scope")
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, val)
		})
	}
}

func TestCreateLabelCommand_ProjectIDFlag(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expected  int64
		expectErr bool
	}{
		{"zero", "0", 0, false},
		{"positive", "123", 123, false},
		{"large", "9999999", 9999999, false},
		{"invalid", "abc", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreateLabelCommand()
			err := cmd.Flags().Set("project", tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				val, _ := cmd.Flags().GetInt64("project")
				assert.Equal(t, tt.expected, val)
			}
		})
	}
}

func TestCreateLabelCommand_GlobalScope_NoProjectIDError(t *testing.T) {
	cmd := CreateLabelCommand()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{})

	assert.NoError(t, cmd.Flags().Set("name", "test-label"))
	assert.NoError(t, cmd.Flags().Set("scope", "g"))
	// Global scope should NOT require project ID - no validation error

	err := cmd.Execute()
	// If there's an error, it should NOT be about project ID (will be API error)
	if err != nil {
		assert.NotContains(t, err.Error(), "project ID is required")
	}
}

func TestCreateLabelCommand_MultipleFlagCombinations(t *testing.T) {
	tests := []struct {
		name        string
		flags       map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "project_scope_without_project_id",
			flags:       map[string]string{"name": "test", "scope": "p"},
			expectError: true,
			errorMsg:    "project ID is required",
		},
		{
			name:        "project_scope_with_zero_project_id",
			flags:       map[string]string{"name": "test", "scope": "p", "project": "0"},
			expectError: true,
			errorMsg:    "project ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreateLabelCommand()

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs([]string{})

			for k, v := range tt.flags {
				assert.NoError(t, cmd.Flags().Set(k, v))
			}

			err := cmd.Execute()
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

func TestCreateLabelCommand_FlagDescriptions(t *testing.T) {
	cmd := CreateLabelCommand()
	flags := cmd.Flags()

	// Verify flag usage descriptions are set
	nameFlag := flags.Lookup("name")
	assert.Contains(t, nameFlag.Usage, "Name of the label")

	colorFlag := flags.Lookup("color")
	assert.Contains(t, colorFlag.Usage, "Color of the label")

	scopeFlag := flags.Lookup("scope")
	assert.Contains(t, scopeFlag.Usage, "Scope of the label")

	projectFlag := flags.Lookup("project")
	assert.Contains(t, projectFlag.Usage, "Id of the project")

	descFlag := flags.Lookup("description")
	assert.Contains(t, descFlag.Usage, "Description of the label")
}
