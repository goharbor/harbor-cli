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

func TestListWorkers_NoWorkers(t *testing.T) {
	out := captureStdout(t, func() {
		ListWorkers([]*models.Worker{})
	})

	assert.Contains(t, out, "No workers found.")
}

func TestListWorkers_WithWorkers(t *testing.T) {
	workers := []*models.Worker{
		{ID: "worker-1", PoolID: "pool-abc123", JobID: "job-9", JobName: "REPLICATE"},
		{ID: "worker-2", PoolID: "pool-abc123", JobID: "job-10", JobName: "RETENTION"},
	}

	out := captureStdout(t, func() {
		ListWorkers(workers)
	})

	assert.Contains(t, out, "WORKER_ID")
	assert.Contains(t, out, "POOL_ID")
	assert.Contains(t, out, "JOB_ID")
	assert.Contains(t, out, "JOB_NAME")
	assert.Contains(t, out, "worker-1")
	assert.Contains(t, out, "job-9")
	assert.Contains(t, out, "REPLICATE")
	assert.Contains(t, out, "Total: 2 worker(s)")
}

func TestListWorkers_NilEntrySkipped(t *testing.T) {
	workers := []*models.Worker{
		nil,
		{ID: "worker-1", PoolID: "pool-abc123", JobID: "job-9", JobName: "REPLICATE"},
	}

	out := captureStdout(t, func() {
		ListWorkers(workers)
	})

	assert.Contains(t, out, "worker-1")
	assert.Contains(t, out, "Total: 2 worker(s)")
}
