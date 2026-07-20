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

	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

func TestHistoryIDColumnWidth(t *testing.T) {
	if got := columns[0]; got.Title != "ID" || got.Width != tablelist.WidthS {
		t.Fatalf("ID column = %#v, want WidthS", got)
	}
}

func TestFormatCreationTimeFallsBackToOriginalTimestamp(t *testing.T) {
	timestamp := "not-a-timestamp"
	if got := formatCreationTime(timestamp); got != timestamp {
		t.Fatalf("formatCreationTime() = %q, want original timestamp %q", got, timestamp)
	}
}

func TestFormatParams(t *testing.T) {
	tests := []struct {
		name   string
		params string
		want   string
	}{
		{name: "empty", params: "", want: "-"},
		{name: "empty object", params: `{}`, want: "-"},
		{name: "invalid JSON", params: "not-json", want: "not-json"},
		{
			name:   "sorted readable values",
			params: `{"workers":2,"delete_untagged":true,"mode":"dry-run","options":{"b":2,"a":1}}`,
			want:   `delete_untagged: true | mode: dry-run | options: {"a":1,"b":2} | workers: 2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatParams(tt.params); got != tt.want {
				t.Fatalf("formatParams() = %q, want %q", got, tt.want)
			}
		})
	}
}
