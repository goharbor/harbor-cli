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
package prompt_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// pointClientAt directs the shared Harbor client at the given base URL for
// the duration of the test, the same way pkg/api tests do.
func pointClientAt(t *testing.T, url string) {
	t.Helper()

	utils.ClientOnce = sync.Once{}
	utils.ClientOnce.Do(func() {})
	utils.ClientInstance = utils.GetClientByConfig(&harbor.ClientSetConfig{URL: url})
	utils.ClientErr = nil

	t.Cleanup(func() {
		utils.ClientOnce = sync.Once{}
		utils.ClientInstance = nil
		utils.ClientErr = nil
	})
}

// A failed API call must surface as an error. It used to be swallowed and
// reported as "no executions", which sent users down the wrong path on a
// plain network failure.
func TestGetRetentionExecutionIDFromUser_PropagatesAPIError(t *testing.T) {
	// unreachable is a port nothing listens on, so dialing it fails
	// immediately without leaving the machine.
	pointClientAt(t, "http://127.0.0.1:1")

	id, err := prompt.GetRetentionExecutionIDFromUser("1")

	require.Error(t, err, "an unreachable server must not read as an empty list")
	assert.Contains(t, err.Error(), "failed to list retention executions")
	assert.Equal(t, int64(0), id)
}

// An empty list is not an error: the caller prints its own message.
func TestGetRetentionExecutionIDFromUser_EmptyList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))
	defer srv.Close()

	pointClientAt(t, srv.URL)

	id, err := prompt.GetRetentionExecutionIDFromUser("1")

	require.NoError(t, err)
	assert.Equal(t, int64(0), id)
}
