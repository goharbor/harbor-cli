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
package schedules

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = writer
	defer func() {
		os.Stdout = originalStdout
	}()

	outputCh := make(chan string, 1)
	go func() {
		var buffer bytes.Buffer
		_, _ = io.Copy(&buffer, reader)
		outputCh <- buffer.String()
	}()

	fn()

	_ = writer.Close()
	return <-outputCh
}

func TestListSchedules(t *testing.T) {
	output := captureOutput(t, func() {
		ListSchedules([]*models.ScheduleTask{
			{
				ID:         42,
				VendorType: "replication",
				VendorID:   7,
				Cron:       "0 0 * * *",
			},
		}, 2, 20, 1)
	})

	for _, want := range []string{"ID", "VENDOR_TYPE", "VENDOR_ID", "CRON", "UPDATE_TIME", "42", "replication", "7", "0 0 * * *", "Page: 2  Page Size: 20  Returned: 1  Total: 1"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got %q", want, output)
		}
	}
}

func TestListSchedulesEmpty(t *testing.T) {
	output := captureOutput(t, func() {
		ListSchedules(nil, 1, 20, 0)
	})

	if !strings.Contains(output, "No schedules found.") {
		t.Fatalf("expected empty-state message, got %q", output)
	}
}

func TestPrintScheduleStatus(t *testing.T) {
	tests := []struct {
		name string
		in   *models.SchedulerStatus
		want string
	}{
		{
			name: "nil",
			in:   nil,
			want: "Scheduler status: unknown",
		},
		{
			name: "paused",
			in:   &models.SchedulerStatus{Paused: true},
			want: "Scheduler status: paused",
		},
		{
			name: "running",
			in:   &models.SchedulerStatus{Paused: false},
			want: "Scheduler status: running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(t, func() {
				PrintScheduleStatus(tt.in)
			})

			if !strings.Contains(output, tt.want) {
				t.Fatalf("expected output to contain %q, got %q", tt.want, output)
			}
		})
	}
}
