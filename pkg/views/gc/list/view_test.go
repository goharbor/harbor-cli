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
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

func TestHistoryIDColumnWidth(t *testing.T) {
	if got := columns[0]; got.Title != "ID" || got.Width != tablelist.WidthS {
		t.Fatalf("ID column = %#v, want WidthS", got)
	}
}

func TestParametersColumnWidth(t *testing.T) {
	if got := columns[4]; got.Title != "Parameters" || got.Width != tablelist.WidthXL*2 {
		t.Fatalf("Parameters column = %#v, want WidthXL * 2", got)
	}
}

func TestHistoryRowsSeparatesSortedParameters(t *testing.T) {
	run := &models.GCHistory{
		ID:            42,
		JobName:       "GARBAGE_COLLECTION",
		JobStatus:     "Success",
		JobKind:       "MANUAL",
		JobParameters: `{"zeta":2,"alpha":1}`,
	}
	want := []table.Row{
		{"42", "GARBAGE_COLLECTION", "Success", "MANUAL", "alpha: 1", "2 hours ago"},
		{"", "", "", "", "zeta: 2", ""},
	}

	if got := historyRows(run, "2 hours ago"); !reflect.DeepEqual(got, want) {
		t.Fatalf("historyRows() = %#v, want %#v", got, want)
	}
}

func TestRenderedHistoryRowsKeepParametersOnSeparateLines(t *testing.T) {
	run := &models.GCHistory{
		ID:            767340,
		JobName:       "GARBAGE_COLLECTION",
		JobStatus:     "Success",
		JobKind:       "MANUAL",
		JobParameters: `{"delete_tag":false,"delete_untagged":false,"dry_run":true,"freed_space":9007199254740991}`,
	}
	rows := historyRows(run, "50 minutes ago")
	rendered := tablelist.NewModel(columns, rows, len(rows)).View()
	t.Logf("rendered GC history table:\n%s", rendered)

	parameters := []string{
		"delete_tag: false",
		"delete_untagged: false",
		"dry_run: true",
		"freed_space: 9007199254740991",
	}
	seenLines := make(map[int]bool, len(parameters))
	lines := strings.Split(rendered, "\n")
	for _, parameter := range parameters {
		lineIndex := -1
		for i, line := range lines {
			if strings.Contains(line, parameter) {
				lineIndex = i
				break
			}
		}
		if lineIndex == -1 {
			t.Fatalf("rendered table does not contain %q", parameter)
		}
		if seenLines[lineIndex] {
			t.Fatalf("parameter %q shares rendered line %d with another parameter", parameter, lineIndex)
		}
		seenLines[lineIndex] = true
	}
}

func TestFormatParams(t *testing.T) {
	tests := []struct {
		name   string
		params string
		want   []string
	}{
		{name: "empty", params: "", want: []string{"-"}},
		{name: "empty object", params: `{}`, want: []string{"-"}},
		{name: "invalid JSON", params: "not-json", want: []string{"not-json"}},
		{
			name:   "sorted values",
			params: `{"workers":2,"delete_untagged":true,"mode":"dry-run","options":{"b":2,"a":1}}`,
			want: []string{
				"delete_untagged: true",
				"mode: dry-run",
				`options: {"a":1,"b":2}`,
				"workers: 2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatParams(tt.params); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("formatParams() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
