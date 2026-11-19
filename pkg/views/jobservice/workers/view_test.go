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

package workers

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

func TestListWorkersPaginationFooter(t *testing.T) {
	output := captureOutput(t, func() {
		ListWorkers([]*models.Worker{
			{ID: "worker-1", JobID: "job-1"},
			{ID: "worker-2"},
		}, 2, 20, 5)
	})

	for _, want := range []string{"worker-1", "worker-2", "Page: 2  Page Size: 20  Returned: 2  Total: 5  Busy: 1"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got %q", want, output)
		}
	}
}

func TestListWorkersEmptyState(t *testing.T) {
	output := captureOutput(t, func() {
		ListWorkers(nil, 1, 20, 0)
	})

	if !strings.Contains(output, "No workers found.") {
		t.Fatalf("expected empty-state message, got %q", output)
	}
}
