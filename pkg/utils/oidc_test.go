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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiscoverOIDCProvider_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/.well-known/openid-configuration", r.URL.Path)
		json.NewEncoder(w).Encode(OIDCDiscovery{
			Issuer:                      "https://idp.example.com",
			DeviceAuthorizationEndpoint: "https://idp.example.com/device/auth",
			TokenEndpoint:               "https://idp.example.com/token",
		})
	}))
	defer server.Close()

	discovery, err := discoverOIDCProvider(server.URL, server.Client())
	assert.NoError(t, err)
	assert.Equal(t, "https://idp.example.com/device/auth", discovery.DeviceAuthorizationEndpoint)
	assert.Equal(t, "https://idp.example.com/token", discovery.TokenEndpoint)
}

func TestDiscoverOIDCProvider_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := discoverOIDCProvider(server.URL, server.Client())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestRequestDeviceAuthorization_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(DeviceAuthorizationResponse{
			DeviceCode:      "device-code-123",
			UserCode:        "ABCD-EFGH",
			VerificationURI: "https://idp.example.com/device",
			ExpiresIn:       300,
			Interval:        5,
		})
	}))
	defer server.Close()

	resp, err := requestDeviceAuthorization(server.URL, "test-client", "openid", server.Client())
	assert.NoError(t, err)
	assert.Equal(t, "ABCD-EFGH", resp.UserCode)
	assert.Equal(t, "device-code-123", resp.DeviceCode)
	assert.Equal(t, 5, resp.Interval)
}

func TestRequestDeviceAuthorization_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DeviceFlowError{
			Error:            "invalid_client",
			ErrorDescription: "Client not found",
		})
	}))
	defer server.Close()

	_, err := requestDeviceAuthorization(server.URL, "bad-client", "openid", server.Client())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid_client")
}

func TestPollForToken_Success(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 3 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(DeviceFlowError{Error: "authorization_pending"})
			return
		}
		json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "access-token-xyz",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			RefreshToken: "refresh-token-xyz",
		})
	}))
	defer server.Close()

	resp, err := pollForToken(server.URL, "device-code", "client-id", 10*time.Millisecond, server.Client())
	assert.NoError(t, err)
	assert.Equal(t, "access-token-xyz", resp.AccessToken)
	assert.Equal(t, "refresh-token-xyz", resp.RefreshToken)
	assert.Equal(t, int64(3600), resp.ExpiresIn)
	assert.Equal(t, 3, callCount)
}

func TestPollForToken_Denied(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(DeviceFlowError{
			Error:            "access_denied",
			ErrorDescription: "User denied the request",
		})
	}))
	defer server.Close()

	_, err := pollForToken(server.URL, "device-code", "client-id", 10*time.Millisecond, server.Client())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access_denied")
}

func TestRunOIDCLogin_Integration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "at-final",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			RefreshToken: "rt-final",
		})
	}))
	defer server.Close()

	originalDiscover := oidcDiscoveryFunc
	originalDeviceAuth := deviceAuthorizationFunc
	originalPoll := pollTokenFunc
	defer func() {
		oidcDiscoveryFunc = originalDiscover
		deviceAuthorizationFunc = originalDeviceAuth
		pollTokenFunc = originalPoll
	}()

	oidcDiscoveryFunc = func(issuerURL string, client *http.Client) (*OIDCDiscovery, error) {
		return &OIDCDiscovery{
			Issuer:                      issuerURL,
			DeviceAuthorizationEndpoint: server.URL + "/device-auth",
			TokenEndpoint:               server.URL + "/token",
		}, nil
	}
	deviceAuthorizationFunc = func(endpoint, clientID, scopes string, client *http.Client) (*DeviceAuthorizationResponse, error) {
		return &DeviceAuthorizationResponse{
			DeviceCode:      "dc-123",
			UserCode:        "XYZW-1234",
			VerificationURI: "https://idp.example.com/device",
			ExpiresIn:       300,
			Interval:        1,
		}, nil
	}
	pollTokenFunc = func(tokenEndpoint, deviceCode, clientID string, interval time.Duration, client *http.Client) (*TokenResponse, error) {
		return pollForToken(tokenEndpoint, deviceCode, clientID, interval, client)
	}

	resp, err := RunOIDCLogin(OIDCFlowOptions{
		IssuerURL: "https://idp.example.com",
		ClientID:  "test-client",
		Scopes:    "openid",
		Client:    server.Client(),
	})
	assert.NoError(t, err)
	assert.Equal(t, "at-final", resp.AccessToken)
	assert.Equal(t, "rt-final", resp.RefreshToken)
}
