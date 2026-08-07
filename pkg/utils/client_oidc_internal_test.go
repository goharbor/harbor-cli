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
package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestOIDCTokenExpiryUnix(t *testing.T) {
	token := testJWTWithExp(time.Now().Add(10 * time.Minute).Unix())

	expiresAt, err := oidcTokenExpiryUnix(token)

	require.NoError(t, err)
	assert.Greater(t, expiresAt, time.Now().Unix())
}

func TestOIDCTokenNeedsRefreshPrefersTokenExpiry(t *testing.T) {
	credential := Credential{
		ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
	}
	token := testJWTWithExp(time.Now().Add(30 * time.Second).Unix())

	assert.True(t, oidcTokenNeedsRefresh(credential, token))
}

func TestOIDCTokenNeedsRefreshFallsBackToStoredExpiry(t *testing.T) {
	credential := Credential{
		ExpiresAt: time.Now().Add(30 * time.Second).Unix(),
	}

	assert.True(t, oidcTokenNeedsRefresh(credential, "not-a-jwt"))
}

func testJWTWithExp(exp int64) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"exp":%d}`, exp)))
	return header + "." + payload + ".signature"
}

func TestOIDCRetryTransportRefreshesAndRetriesOnUnauthorized(t *testing.T) {
	requests := 0
	baseTransport := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		requests++
		resp := &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    r,
		}
		switch r.Header.Get("Authorization") {
		case "Bearer stale-token":
			resp.StatusCode = http.StatusUnauthorized
		case "Bearer fresh-token":
			resp.StatusCode = http.StatusOK
		}
		return resp, nil
	})

	tokenManager := &oidcTokenManager{
		token: "stale-token",
		refreshFn: func(_ Credential) (string, error) {
			return "fresh-token", nil
		},
	}

	transport := &oidcRetryTransport{
		base:         baseTransport,
		tokenManager: tokenManager,
	}

	req, err := http.NewRequest(http.MethodGet, "https://harbor.example.com/api/v2.0/projects", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenManager.Token())

	resp, err := transport.RoundTrip(req)
	defer resp.Body.Close()
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "fresh-token", tokenManager.Token())
	assert.Equal(t, 2, requests)
}

func TestCloneRequestForRetryWithReplayableBody(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "https://example.com", strings.NewReader("payload"))
	require.NoError(t, err)

	cloned, err := cloneRequestForRetry(req)

	require.NoError(t, err)
	body, err := io.ReadAll(cloned.Body)
	require.NoError(t, err)
	assert.Equal(t, "payload", string(body))
}
