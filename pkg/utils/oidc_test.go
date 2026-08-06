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
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func withDefaultTransport(t *testing.T, transport http.RoundTripper) {
	t.Helper()
	original := http.DefaultTransport
	http.DefaultTransport = transport
	t.Cleanup(func() {
		http.DefaultTransport = original
	})
}

func TestInitiateOIDCLogin(t *testing.T) {
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/c/oidc/login", r.URL.Path)
		assert.Equal(t, "cli", r.URL.Query().Get("mode"))
		body, err := json.Marshal(utils.OIDCLoginResponse{
			RedirectURL: "https://idp.example/authorize",
			PollToken:   "poll-token-1",
		})
		require.NoError(t, err)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Request:    r,
		}, nil
	}))

	resp, err := utils.InitiateOIDCLogin("https://harbor.example.com")

	require.NoError(t, err)
	assert.Equal(t, "https://idp.example/authorize", resp.RedirectURL)
	assert.Equal(t, "poll-token-1", resp.PollToken)
}

func TestInitiateOIDCLoginPreservesBasePath(t *testing.T) {
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/harbor/c/oidc/login", r.URL.Path)
		assert.Equal(t, "cli", r.URL.Query().Get("mode"))
		body, err := json.Marshal(utils.OIDCLoginResponse{
			RedirectURL: "https://idp.example/authorize",
			PollToken:   "poll-token-1",
		})
		require.NoError(t, err)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Request:    r,
		}, nil
	}))

	resp, err := utils.InitiateOIDCLogin("https://harbor.example.com/harbor")

	require.NoError(t, err)
	assert.Equal(t, "https://idp.example/authorize", resp.RedirectURL)
	assert.Equal(t, "poll-token-1", resp.PollToken)
}

func TestInitiateOIDCLoginRejectsBrowserRedirectResponse(t *testing.T) {
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusFound,
			Header:     http.Header{"Location": []string{"https://accounts.example.com/authorize"}},
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    r,
		}, nil
	}))

	resp, err := utils.InitiateOIDCLogin("https://harbor.example.com")

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "may not support CLI OIDC login yet")
}

func TestInitiateOIDCLoginRejectsHTMLResponse(t *testing.T) {
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/html; charset=utf-8"}},
			Body:       io.NopCloser(strings.NewReader(`<a href="https://accounts.example.com/authorize">Found</a>.`)),
			Request:    r,
		}, nil
	}))

	resp, err := utils.InitiateOIDCLogin("https://harbor.example.com")

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "may not support CLI OIDC login yet")
	assert.ErrorContains(t, err, "text/html")
}

func TestPollForOIDCTokenReady(t *testing.T) {
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, "/c/oidc/cli-token", r.URL.Path)
		assert.Equal(t, "poll-token-1", r.URL.Query().Get("poll_token"))
		body, err := json.Marshal(utils.OIDCPollResponse{
			Status:   "ready",
			IDToken:  "id-token",
			Username: "alice",
		})
		require.NoError(t, err)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Request:    r,
		}, nil
	}))

	resp, err := utils.PollForOIDCToken("https://harbor.example.com", "poll-token-1", time.Second)

	require.NoError(t, err)
	assert.Equal(t, "ready", resp.Status)
	assert.Equal(t, "id-token", resp.IDToken)
	assert.Equal(t, "alice", resp.Username)
}

func TestPollForOIDCTokenFailed(t *testing.T) {
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		body, err := json.Marshal(utils.OIDCPollResponse{
			Status: "failed",
			Error:  "state expired",
		})
		require.NoError(t, err)
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Request:    r,
		}, nil
	}))

	resp, err := utils.PollForOIDCToken("https://harbor.example.com", "poll-token-1", time.Second)

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "state expired")
}

func TestPollForOIDCTokenExpired(t *testing.T) {
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		body, err := json.Marshal(utils.OIDCPollResponse{
			Status: "expired",
		})
		require.NoError(t, err)
		return &http.Response{
			StatusCode: http.StatusGone,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Request:    r,
		}, nil
	}))

	resp, err := utils.PollForOIDCToken("https://harbor.example.com", "poll-token-1", time.Second)

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "expired before token retrieval")
}

func TestPollForOIDCTokenTimeoutWhilePending(t *testing.T) {
	var requests int32
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt32(&requests, 1)
		assert.Equal(t, "/c/oidc/cli-token", r.URL.Path)
		assert.Equal(t, "poll-token-1", r.URL.Query().Get("poll_token"))
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    r,
		}, nil
	}))

	resp, err := utils.PollForOIDCToken("https://harbor.example.com", "poll-token-1", 100*time.Millisecond)

	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "timed out waiting for OIDC authentication")
	assert.Equal(t, int32(1), atomic.LoadInt32(&requests))
}

func TestRefreshOIDCToken(t *testing.T) {
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/c/oidc/refresh", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req utils.OIDCRefreshRequest
		require.NoError(t, json.Unmarshal(body, &req))
		assert.Equal(t, "refresh-token-1", req.RefreshToken)

		respBody, err := json.Marshal(utils.OIDCRefreshResponse{
			IDToken:      "id-token-2",
			RefreshToken: "refresh-token-2",
			ExpiresAt:    1234567890,
		})
		require.NoError(t, err)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(string(respBody))),
			Request:    r,
		}, nil
	}))

	resp, err := utils.RefreshOIDCToken("https://harbor.example.com", "refresh-token-1")

	require.NoError(t, err)
	assert.Equal(t, "id-token-2", resp.IDToken)
	assert.Equal(t, "refresh-token-2", resp.RefreshToken)
	assert.Equal(t, int64(1234567890), resp.ExpiresAt)
}
