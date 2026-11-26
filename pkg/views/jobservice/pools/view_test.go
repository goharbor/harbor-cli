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
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = oldStdout
	_, _ = io.Copy(&buf, r)

	return buf.String()
}

func TestListPools_NoPools(t *testing.T) {
	out := captureStdout(t, func() {
		ListPools([]*models.WorkerPool{})
	})

	assert.Contains(t, out, "No worker pools found.")
}

func TestListPools_WithPools(t *testing.T) {
	pools := []*models.WorkerPool{
		{WorkerPoolID: "pool-abc123", Host: "10.0.0.1:8080", Pid: 12345, Concurrency: 10},
		{WorkerPoolID: "pool-def456", Host: "10.0.0.2:8080", Pid: 54321, Concurrency: 20},
	}

	out := captureStdout(t, func() {
		ListPools(pools)
	})

	assert.Contains(t, out, "POOL_ID")
	assert.Contains(t, out, "HOST")
	assert.Contains(t, out, "PID")
	assert.Contains(t, out, "CONCURRENCY")
	assert.Contains(t, out, "START_AT")
	assert.Contains(t, out, "pool-abc123")
	assert.Contains(t, out, "10.0.0.1:8080")
	assert.Contains(t, out, "pool-def456")
	assert.Contains(t, out, "Total: 2 pool(s)")
}

func TestListPools_NilEntrySkipped(t *testing.T) {
	pools := []*models.WorkerPool{
		nil,
		{WorkerPoolID: "pool-abc123", Host: "10.0.0.1:8080", Pid: 12345, Concurrency: 10},
	}

	out := captureStdout(t, func() {
		ListPools(pools)
	})

	assert.Contains(t, out, "pool-abc123")
	assert.Contains(t, out, "Total: 2 pool(s)")
}
