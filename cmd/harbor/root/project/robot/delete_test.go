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
package robot

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

// useEmptyRobotServer answers robot listings with an empty result, so
// lookups by name find nothing.
func useEmptyRobotServer(t *testing.T) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2.0/configurations":
			_, _ = w.Write([]byte(`{"robot_name_prefix": {"value": "robot$", "editable": true}}`))
		case "/api/v2.0/projects/library":
			_, _ = w.Write([]byte(`{"project_id": 1, "name": "library"}`))
		case "/api/v2.0/robots":
			_, _ = w.Write([]byte("[]"))
		default:
			http.Error(w, "unexpected request: "+r.URL.Path, http.StatusNotFound)
		}
	}))

	utils.ClientOnce = sync.Once{}
	utils.ClientOnce.Do(func() {})
	utils.ClientInstance = utils.GetClientByConfig(&harbor.ClientSetConfig{URL: srv.URL})
	utils.ClientErr = nil

	t.Cleanup(func() {
		srv.Close()
		utils.ClientOnce = sync.Once{}
		utils.ClientInstance = nil
		utils.ClientErr = nil
	})
}

// Regression: a failed lookup by name used to be printed and then
// dereferenced, crashing the CLI instead of returning the error.
func TestDeleteRobotCommandReturnsLookupError(t *testing.T) {
	useEmptyRobotServer(t)

	cmd := DeleteRobotCommand()
	cmd.SetArgs([]string{"robot$library+missing", "--project", "library"})
	err := cmd.Execute()
	require.ErrorContains(t, err, "failed to get robot by name")
}
