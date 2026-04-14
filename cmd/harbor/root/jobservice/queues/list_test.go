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

package queues

import (
	"reflect"
	"testing"
)

func TestNormalizeJobTypes(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "deduplicates and trims",
			in:   []string{"REPLICATION", " REPLICATION ", "RETENTION"},
			want: []string{"REPLICATION", "RETENTION"},
		},
		{
			name: "supports comma separated",
			in:   []string{"REPLICATION,RETENTION", "GC"},
			want: []string{"REPLICATION", "RETENTION", "GC"},
		},
		{
			name: "all short circuits",
			in:   []string{"REPLICATION", "all", "RETENTION"},
			want: []string{"all"},
		},
		{
			name: "case insensitive dedupe",
			in:   []string{"Replication", "replication", "RETENTION"},
			want: []string{"Replication", "RETENTION"},
		},
		{
			name: "empty values removed",
			in:   []string{"", "  ", ",", "REPLICATION"},
			want: []string{"REPLICATION"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeJobTypes(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("normalizeJobTypes() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestShouldIncludeQueueForAction(t *testing.T) {
	tests := []struct {
		name   string
		action string
		paused bool
		want   bool
	}{
		{name: "resume paused", action: "resume", paused: true, want: true},
		{name: "resume unpaused", action: "resume", paused: false, want: false},
		{name: "pause paused", action: "pause", paused: true, want: false},
		{name: "pause unpaused", action: "pause", paused: false, want: true},
		{name: "stop paused", action: "stop", paused: true, want: true},
		{name: "stop unpaused", action: "stop", paused: false, want: true},
		{name: "unknown action", action: "unknown", paused: false, want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldIncludeQueueForAction(tc.action, tc.paused)
			if got != tc.want {
				t.Fatalf("shouldIncludeQueueForAction(%q, %v) = %v, want %v", tc.action, tc.paused, got, tc.want)
			}
		})
	}
}

func TestActionLabel(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: "Updating"},
		{name: "lower", in: "pause", want: "Pause"},
		{name: "upper", in: "RESUME", want: "Resume"},
		{name: "mixed", in: "sToP", want: "Stop"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := actionLabel(tc.in)
			if got != tc.want {
				t.Fatalf("actionLabel(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
