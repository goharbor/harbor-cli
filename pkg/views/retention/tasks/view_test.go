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
package tasks

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		maxLen int
		want   string
	}{
		{name: "short string untouched", s: "abc", maxLen: 10, want: "abc"},
		{name: "exact length untouched", s: "abcdef", maxLen: 6, want: "abcdef"},
		{name: "long string ellipsized", s: "abcdefghij", maxLen: 6, want: "abc..."},
		{name: "width below 3 does not panic", s: "abcdef", maxLen: 2, want: "ab"},
		{name: "zero width returns empty", s: "abcdef", maxLen: 0, want: ""},
		{name: "negative width returns empty", s: "abcdef", maxLen: -5, want: ""},
		{name: "multibyte not split mid-rune", s: "日本語のテキスト", maxLen: 5, want: "日本..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, truncateString(tt.s, tt.maxLen))
		})
	}
}

// The success path must render every task field into its own column.
func TestBuildTaskRows(t *testing.T) {
	tasks := []*models.RetentionExecutionTask{
		{
			ID:             11,
			ExecutionID:    7,
			Repository:     "library/nginx",
			Status:         "Success",
			Retained:       3,
			Total:          10,
			StatusCode:     200,
			StatusRevision: 42,
			StartTime:      "2026-11-10 10:00:00",
			EndTime:        "2026-11-10 10:00:05",
		},
	}

	// Fixed wide columns so the expectation does not depend on the
	// terminal width of the machine running the test.
	columns := make([]table.Column, 10)
	for i := range columns {
		columns[i] = table.Column{Title: "col", Width: 40}
	}
	rows := buildTaskRows(tasks, columns)

	require.Len(t, rows, 1)
	require.Len(t, rows[0], len(columns), "one cell per column")

	expected := []string{
		"11", "7", "library/nginx", "Success", "3", "10", "200", "42",
		"2026-11-10 10:00:00", "2026-11-10 10:00:05",
	}
	for i, want := range expected {
		assert.Equal(t, want, rows[0][i], "column %d (%s)", i, columns[i].Title)
	}
}

// Cells wider than their column must be truncated, not overflow it.
func TestBuildTaskRowsTruncatesToColumnWidth(t *testing.T) {
	tasks := []*models.RetentionExecutionTask{
		{ID: 1, Repository: strings.Repeat("r", 50)},
	}

	columns := getAdjustedColumns()
	rows := buildTaskRows(tasks, columns)

	require.Len(t, rows, 1)
	for i, cell := range rows[0] {
		assert.LessOrEqual(t, len([]rune(cell)), columns[i].Width,
			"column %s must not overflow", columns[i].Title)
	}
}
