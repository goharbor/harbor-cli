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
package replication

import (
	"strings"
	"testing"
)

func TestLogsCommandAcceptsAdvertisedPositionalArgs(t *testing.T) {
	cmd := LogsCommand()

	if err := cmd.Args(cmd, []string{"123", "456"}); err != nil {
		t.Fatalf("expected two positional args to be accepted, got %v", err)
	}

	if err := cmd.Args(cmd, []string{"123", "456", "789"}); err == nil {
		t.Fatal("expected more than two positional args to be rejected")
	}
}

func TestApplyLogArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		execID    int64
		taskID    int64
		wantExec  int64
		wantTask  int64
		wantError string
	}{
		{
			name:     "execution and task IDs",
			args:     []string{"123", "456"},
			wantExec: 123,
			wantTask: 456,
		},
		{
			name:     "execution ID only",
			args:     []string{"123"},
			wantExec: 123,
		},
		{
			name:      "invalid execution ID",
			args:      []string{"abc"},
			wantError: "invalid replication execution ID",
		},
		{
			name:      "invalid task ID",
			args:      []string{"123", "abc"},
			wantError: "invalid replication task ID",
		},
		{
			name:      "execution ID flag and argument conflict",
			args:      []string{"123"},
			execID:    99,
			wantExec:  99,
			wantError: "execution ID cannot be provided both as a flag and an argument",
		},
		{
			name:      "task ID flag and argument conflict",
			args:      []string{"123", "456"},
			taskID:    99,
			wantExec:  123,
			wantTask:  99,
			wantError: "task ID cannot be provided both as a flag and an argument",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			execID := tc.execID
			taskID := tc.taskID

			err := applyLogArgs(tc.args, &execID, &taskID)
			if tc.wantError != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", tc.wantError)
				}
				if !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("expected error containing %q, got %v", tc.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if execID != tc.wantExec {
				t.Fatalf("expected execution ID %d, got %d", tc.wantExec, execID)
			}
			if taskID != tc.wantTask {
				t.Fatalf("expected task ID %d, got %d", tc.wantTask, taskID)
			}
		})
	}
}
