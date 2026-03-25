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

import (
	"regexp"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

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
		expectedRx   string // regex pattern to match (for times that vary)
		wantErr      bool
	}{
		{
			name:       "returns base query only",
			baseQuery:  "operation=push",
			expectedRx: "^operation=push$",
		},
		{
			name:         "builds query with convenience filters",
			baseQuery:    "operation_result=true",
			operation:    "create_artifact",
			resourceType: "artifact",
			resource:     "library/nginx",
			username:     "admin",
			expectedRx:   "^operation_result=true,operation=create_artifact,resource_type=artifact,resource=library/nginx,username=admin$",
		},
		{
			name:       "builds range query with both times specified",
			fromTime:   "2025-01-01T01:02:03Z",
			toTime:     "2025-01-01 05:06:07",
			expectedRx: "^op_time=\\[2025-01-01 01:02:03~2025-01-01 05:06:07\\]$",
		},
		{
			name:       "from-time alone defaults to-time to current time",
			fromTime:   "2025-01-01T01:02:03Z",
			expectedRx: "^op_time=\\[2025-01-01 01:02:03~.*\\]$", // matches any end time
		},
		{
			name:    "to-time alone is rejected",
			toTime:  "2025-01-01 05:06:07",
			wantErr: true,
		},
		{
			name:     "fails for invalid from time",
			fromTime: "invalid-time",
			toTime:   "2025-01-01 05:06:07",
			wantErr:  true,
		},
		{
			name:     "fails for invalid to time",
			fromTime: "2025-01-01T01:02:03Z",
			toTime:   "invalid-time",
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

			matched, _ := regexp.MatchString(tt.expectedRx, query)
			if !matched {
				t.Fatalf("expected query to match regex %q, got %q", tt.expectedRx, query)
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

func TestPaginateAuditLogEventTypes(t *testing.T) {
	eventTypes := []*models.AuditLogEventType{
		{EventType: "event-1"},
		{EventType: "event-2"},
		{EventType: "event-3"},
		{EventType: "event-4"},
	}

	tests := []struct {
		name      string
		page      int64
		pageSize  int64
		wantLen   int
		firstItem string
		wantErr   bool
	}{
		{
			name:      "first page",
			page:      1,
			pageSize:  2,
			wantLen:   2,
			firstItem: "event-1",
		},
		{
			name:      "second page",
			page:      2,
			pageSize:  2,
			wantLen:   2,
			firstItem: "event-3",
		},
		{
			name:     "page out of range",
			page:     3,
			pageSize: 2,
			wantLen:  0,
		},
		{
			name:     "invalid page",
			page:     0,
			pageSize: 2,
			wantErr:  true,
		},
		{
			name:     "invalid page size",
			page:     1,
			pageSize: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := paginateAuditLogEventTypes(eventTypes, tt.page, tt.pageSize)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != tt.wantLen {
				t.Fatalf("expected len %d, got %d", tt.wantLen, len(got))
			}

			if tt.firstItem != "" && got[0].EventType != tt.firstItem {
				t.Fatalf("expected first item %q, got %q", tt.firstItem, got[0].EventType)
			}
		})
	}
}

func TestAuditLogEventTypeName(t *testing.T) {
	tests := []struct {
		name      string
		eventType *models.AuditLogEventType
		expected  string
	}{
		{
			name:      "returns event type name",
			eventType: &models.AuditLogEventType{EventType: "create_artifact"},
			expected:  "create_artifact",
		},
		{
			name:      "trims whitespace",
			eventType: &models.AuditLogEventType{EventType: "  login_user  "},
			expected:  "login_user",
		},
		{
			name:      "nil event type",
			eventType: nil,
			expected:  "-",
		},
		{
			name:      "empty event type",
			eventType: &models.AuditLogEventType{EventType: "   "},
			expected:  "-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := auditLogEventTypeName(tt.eventType)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
