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

//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// projectStorageLimit returns the storage quota Harbor holds for project.
func projectStorageLimit(t *testing.T, project string) int64 {
	t.Helper()

	var p struct {
		ProjectID int64 `json:"project_id"`
	}
	doJSON(t, http.MethodGet, "/api/v2.0/projects/"+url.PathEscape(project), nil, &p)

	var quotas []struct {
		Hard map[string]int64 `json:"hard"`
	}
	doJSON(t, http.MethodGet, fmt.Sprintf("/api/v2.0/quotas?reference=project&reference_id=%d", p.ProjectID), nil, &quotas)
	require.Len(t, quotas, 1)

	return quotas[0].Hard["storage"]
}

// A storage limit with a unit must reach Harbor as bytes, not be dropped.
func Test_ProjectCreate_StorageLimitWithUnit(t *testing.T) {
	loginAndReload(t)

	project := fmt.Sprintf("e2e-873-%d", time.Now().UnixNano())
	deleteProjectOnCleanup(t, project)

	require.NoError(t, harbor(t, "project", "create", project, "--storage-limit", "10GiB"))
	require.Equal(t, int64(10<<30), projectStorageLimit(t, project))
}

// An invalid storage limit must fail instead of creating the project.
func Test_ProjectCreate_InvalidStorageLimit(t *testing.T) {
	loginAndReload(t)

	project := fmt.Sprintf("e2e-873-%d", time.Now().UnixNano())
	deleteProjectOnCleanup(t, project)

	require.Error(t, harbor(t, "project", "create", project, "--storage-limit", "10GB"))
	require.Equal(t, http.StatusNotFound, projectStatus(t, http.MethodGet, project))
}
