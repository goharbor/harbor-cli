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
	"errors"
	"strings"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

func TestListCommandAddsPaginationFlags(t *testing.T) {
	cmd := ListCommand()

	if got := cmd.Flags().Lookup("page").DefValue; got != "1" {
		t.Fatalf("expected default page to be 1, got %s", got)
	}

	if got := cmd.Flags().Lookup("page-size").DefValue; got != "20" {
		t.Fatalf("expected default page-size to be 20, got %s", got)
	}
}

func TestListWorkersPaginationValidation(t *testing.T) {
	if err := listWorkers(ListCommand(), nil, "", false, 0, 20); err == nil || !strings.Contains(err.Error(), "page must be >= 1") {
		t.Fatalf("expected page validation error, got %v", err)
	}

	if err := listWorkers(ListCommand(), nil, "", false, 1, 0); err == nil || !strings.Contains(err.Error(), "page-size must be >= 1") {
		t.Fatalf("expected page-size validation error, got %v", err)
	}
}

func TestListWorkersPageWindow(t *testing.T) {
	workers := []*models.Worker{
		{ID: "worker-1"},
		{ID: "worker-2"},
		{ID: "worker-3"},
	}

	page := int64(2)
	pageSize := int64(1)
	start := int((page - 1) * pageSize)
	end := int(page * pageSize)
	pageWorkers := workers[start:end]

	if len(pageWorkers) != 1 || pageWorkers[0].ID != "worker-2" {
		t.Fatalf("expected second page to contain worker-2, got %#v", pageWorkers)
	}
}

func TestFormatWorkerActionError(t *testing.T) {
	got := formatWorkerActionError("failed to free worker", errors.New("[POST /jobservice/jobs][404] (status 404)"))
	if got == nil || !strings.Contains(got.Error(), "job not found or already completed") {
		t.Fatalf("expected 404 mapping, got %v", got)
	}
}
