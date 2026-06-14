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
package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitiateOIDCLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/c/oidc/cli/login", r.URL.Path)
		_ = json.NewEncoder(w).Encode(utils.OIDCLoginResponse{
			RedirectURL:    "https://idp.example/authorize",
			TransactionID: "tx-1",
		})
	}))
	defer server.Close()

	resp, err := utils.InitiateOIDCLogin(server.URL)

	require.NoError(t, err)
	assert.Equal(t, "https://idp.example/authorize", resp.RedirectURL)
	assert.Equal(t, "tx-1", resp.TransactionID)
}

func TestPollForOIDCTokenReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/c/oidc/cli-token", r.URL.Path)
		assert.Equal(t, "tx-1", r.URL.Query().Get("transaction_id"))
		_ = json.NewEncoder(w).Encode(utils.OIDCPollResponse{
			Status:   "ready",
			IDToken:  "id-token",
			Username: "alice",
		})
	}))
	defer server.Close()

	resp, err := utils.PollForOIDCToken(server.URL, "tx-1", time.Second)

	require.NoError(t, err)
	assert.Equal(t, "ready", resp.Status)
	assert.Equal(t, "id-token", resp.IDToken)
	assert.Equal(t, "alice", resp.Username)
}

func TestPollForOIDCTokenFailed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(utils.OIDCPollResponse{
			Status: "failed",
			Error:  "state expired",
		})
	}))
	defer server.Close()

	resp, err := utils.PollForOIDCToken(server.URL, "tx-1", time.Second)

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "state expired")
}
