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

package list

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "string shorter than maxLen",
			input:    "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "string equal to maxLen",
			input:    "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "string longer than maxLen",
			input:    "hello world",
			maxLen:   8,
			expected: "hello...",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   5,
			expected: "",
		},
		{
			name:     "maxLen less than 3",
			input:    "hello",
			maxLen:   2,
			expected: "he",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestFormatScopeSelectors(t *testing.T) {
	tests := []struct {
		name             string
		input            map[string][]models.RetentionSelector
		expected         string
		expectedContains []string
	}{
		{
			name:     "empty map",
			input:    map[string][]models.RetentionSelector{},
			expected: "",
		},
		{
			name: "single key with one selector",
			input: map[string][]models.RetentionSelector{
				"repository": {
					{Decoration: "matches", Pattern: "repo*"},
				},
			},
			expected: "repository: [{matches   repo*}]",
		},
		{
			name: "single key with multiple selectors",
			input: map[string][]models.RetentionSelector{
				"repository": {
					{Decoration: "matches", Pattern: "repo*"},
					{Decoration: "excludes", Pattern: "temp*"},
				},
			},
			expected: "repository: [{matches   repo*}, {excludes   temp*}]",
		},
		{
			name: "multiple keys",
			input: map[string][]models.RetentionSelector{
				"repository": {
					{Decoration: "matches", Pattern: "repo*"},
				},
				"other": {
					{Decoration: "excludes", Pattern: "test*"},
				},
			},
			expectedContains: []string{
				"repository: [{matches   repo*}]",
				"other: [{excludes   test*}]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatScopeSelectors(tt.input)
			if tt.expected != "" || tt.expectedContains == nil {
				if result != tt.expected {
					t.Errorf("formatScopeSelectors() = %q, want %q", result, tt.expected)
				}
			}
			if tt.expectedContains != nil {
				for _, substr := range tt.expectedContains {
					if !containsString(result, substr) {
						t.Errorf("formatScopeSelectors() = %q, expected to contain %q", result, substr)
					}
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}