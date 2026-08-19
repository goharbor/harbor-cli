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

package pools

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
	defer func() { os.Stdout = originalStdout }()

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

func TestListPoolsEmptyState(t *testing.T) {
	output := captureOutput(t, func() {
		ListPools(nil)
	})

	if !strings.Contains(output, "No worker pools found.") {
		t.Fatalf("expected empty-state message, got %q", output)
	}
}

func TestListPoolsRendersRows(t *testing.T) {
	output := captureOutput(t, func() {
		ListPools([]*models.WorkerPool{
			{
				WorkerPoolID: "pool-1",
				Pid:          123,
				Concurrency:  5,
				Host:         "worker-host",
			},
		})
	})

	for _, want := range []string{"POOL_ID", "pool-1", "123", "5", "worker-host", "Total: 1 worker pool(s)"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got %q", want, output)
		}
	}
}
