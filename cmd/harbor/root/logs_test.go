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
package root

import "testing"

func TestBuildAuditLogQuery(t *testing.T) {
	tests := []struct {
		name         string
		baseQuery    string
		operation    string
		resourceType string
		resource     string
		username     string
		fromTime     string
		toTime       string
		expected     string
		wantErr      bool
	}{
		{
			name:      "returns base query only",
			baseQuery: "operation=push",
			expected:  "operation=push",
		},
		{
			name:         "builds query with convenience filters",
			baseQuery:    "operation_result=true",
			operation:    "create_artifact",
			resourceType: "artifact",
			resource:     "library/nginx",
			username:     "admin",
			expected:     "operation_result=true,operation=create_artifact,resource_type=artifact,resource=library/nginx,username=admin",
		},
		{
			name:     "builds range query with normalized times",
			fromTime: "2025-01-01T01:02:03Z",
			toTime:   "2025-01-01 05:06:07",
			expected: "op_time=[2025-01-01 01:02:03~2025-01-01 05:06:07]",
		},
		{
			name:    "fails when one range bound is missing",
			toTime:  "2025-01-01 05:06:07",
			wantErr: true,
		},
		{
			name:     "fails for invalid from time",
			fromTime: "invalid-time",
			toTime:   "2025-01-01 05:06:07",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := buildAuditLogQuery(
				tt.baseQuery,
				tt.operation,
				tt.resourceType,
				tt.resource,
				tt.username,
				tt.fromTime,
				tt.toTime,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if query != tt.expected {
				t.Fatalf("expected query %q, got %q", tt.expected, query)
			}
		})
	}
}

func TestNormalizeAuditTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "accepts RFC3339",
			input:    "2025-01-01T01:02:03Z",
			expected: "2025-01-01 01:02:03",
		},
		{
			name:     "accepts plain datetime",
			input:    "2025-01-01 01:02:03",
			expected: "2025-01-01 01:02:03",
		},
		{
			name:    "fails for invalid input",
			input:   "2025/01/01",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeAuditTime(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.expected {
				t.Fatalf("expected normalized time %q, got %q", tt.expected, got)
			}
		})
	}
}
