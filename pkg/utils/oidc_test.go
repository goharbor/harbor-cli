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
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitiateOIDCLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/c/oidc/login", r.URL.Path)
		assert.Equal(t, "cli", r.URL.Query().Get("mode"))
		_ = json.NewEncoder(w).Encode(utils.OIDCLoginResponse{
			RedirectURL: "https://idp.example/authorize",
			PollToken:   "poll-token-1",
		})
	}))
	defer server.Close()

	resp, err := utils.InitiateOIDCLogin(server.URL)

	require.NoError(t, err)
	assert.Equal(t, "https://idp.example/authorize", resp.RedirectURL)
	assert.Equal(t, "poll-token-1", resp.PollToken)
}

func TestInitiateOIDCLoginPreservesBasePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/harbor/c/oidc/login", r.URL.Path)
		assert.Equal(t, "cli", r.URL.Query().Get("mode"))
		_ = json.NewEncoder(w).Encode(utils.OIDCLoginResponse{
			RedirectURL: "https://idp.example/authorize",
			PollToken:   "poll-token-1",
		})
	}))
	defer server.Close()

	resp, err := utils.InitiateOIDCLogin(server.URL + "/harbor")

	require.NoError(t, err)
	assert.Equal(t, "https://idp.example/authorize", resp.RedirectURL)
	assert.Equal(t, "poll-token-1", resp.PollToken)
}

func TestPollForOIDCTokenReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/c/oidc/cli-token", r.URL.Path)
		assert.Equal(t, "poll-token-1", r.URL.Query().Get("poll_token"))
		_ = json.NewEncoder(w).Encode(utils.OIDCPollResponse{
			Status:   "ready",
			IDToken:  "id-token",
			Username: "alice",
		})
	}))
	defer server.Close()

	resp, err := utils.PollForOIDCToken(server.URL, "poll-token-1", time.Second)

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

	resp, err := utils.PollForOIDCToken(server.URL, "poll-token-1", time.Second)

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "state expired")
}

func TestPollForOIDCTokenExpired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGone)
		_ = json.NewEncoder(w).Encode(utils.OIDCPollResponse{
			Status: "expired",
		})
	}))
	defer server.Close()

	resp, err := utils.PollForOIDCToken(server.URL, "poll-token-1", time.Second)

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "expired before token retrieval")
}

func TestPollForOIDCTokenTimeoutWhilePending(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		assert.Equal(t, "/c/oidc/cli-token", r.URL.Path)
		assert.Equal(t, "poll-token-1", r.URL.Query().Get("poll_token"))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	resp, err := utils.PollForOIDCToken(server.URL, "poll-token-1", 100*time.Millisecond)

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "timed out waiting for OIDC authentication")
	assert.Equal(t, int32(1), atomic.LoadInt32(&requests))
}

func TestRefreshOIDCToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/c/oidc/refresh", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req utils.OIDCRefreshRequest
		require.NoError(t, json.Unmarshal(body, &req))
		assert.Equal(t, "refresh-token-1", req.RefreshToken)

		_ = json.NewEncoder(w).Encode(utils.OIDCRefreshResponse{
			IDToken:      "id-token-2",
			RefreshToken: "refresh-token-2",
			ExpiresAt:    1234567890,
		})
	}))
	defer server.Close()

	resp, err := utils.RefreshOIDCToken(server.URL, "refresh-token-1")

	require.NoError(t, err)
	assert.Equal(t, "id-token-2", resp.IDToken)
	assert.Equal(t, "refresh-token-2", resp.RefreshToken)
	assert.Equal(t, int64(1234567890), resp.ExpiresAt)
}
