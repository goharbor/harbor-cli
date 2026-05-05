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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultOIDCScopes   = "openid profile email offline_access"
	defaultPollInterval = 5 * time.Second
)

type OIDCDiscovery struct {
	Issuer                      string   `json:"issuer"`
	DeviceAuthorizationEndpoint string   `json:"device_authorization_endpoint"`
	TokenEndpoint               string   `json:"token_endpoint"`
	AuthorizationEndpoint       string   `json:"authorization_endpoint"`
	ScopesSupported             []string `json:"scopes_supported"`
	GrantTypesSupported         []string `json:"grant_types_supported"`
}

type DeviceAuthorizationResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

type DeviceFlowError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type OIDCFlowOptions struct {
	IssuerURL string
	ClientID  string
	Scopes    string
	Client    *http.Client
}

var oidcDiscoveryFunc = discoverOIDCProvider

func discoverOIDCProvider(issuerURL string, client *http.Client) (*OIDCDiscovery, error) {
	discoveryURL := strings.TrimRight(issuerURL, "/") + "/.well-known/openid-configuration"
	resp, err := client.Get(discoveryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover OIDC provider at %s: %w", discoveryURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OIDC discovery failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var discovery OIDCDiscovery
	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		return nil, fmt.Errorf("failed to parse OIDC discovery response: %w", err)
	}
	return &discovery, nil
}

var deviceAuthorizationFunc = requestDeviceAuthorization

func requestDeviceAuthorization(endpoint, clientID, scopes string, client *http.Client) (*DeviceAuthorizationResponse, error) {
	data := url.Values{
		"client_id": {clientID},
	}
	if scopes != "" {
		data.Set("scope", scopes)
	}

	resp, err := client.PostForm(endpoint, data)
	if err != nil {
		return nil, fmt.Errorf("device authorization request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		deviceErr, err := parseDeviceFlowError(resp)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("device authorization failed: %s - %s", deviceErr.Error, deviceErr.ErrorDescription)
	}

	var deviceResp DeviceAuthorizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return nil, fmt.Errorf("failed to parse device authorization response: %w", err)
	}
	return &deviceResp, nil
}

var pollTokenFunc = pollForToken

func pollForToken(tokenEndpoint, deviceCode, clientID string, interval time.Duration, client *http.Client) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"device_code": {deviceCode},
		"client_id":   {clientID},
	}

	for {
		resp, err := client.PostForm(tokenEndpoint, data)
		if err != nil {
			return nil, fmt.Errorf("token request failed: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			var tokenResp TokenResponse
			if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
				resp.Body.Close()
				return nil, fmt.Errorf("failed to parse token response: %w", err)
			}
			resp.Body.Close()
			return &tokenResp, nil
		}

		deviceErr, err := parseDeviceFlowError(resp)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if deviceErr.Error != "authorization_pending" && deviceErr.Error != "slow_down" {
			return nil, fmt.Errorf("device authorization failed: %s - %s", deviceErr.Error, deviceErr.ErrorDescription)
		}

		if deviceErr.Error == "slow_down" {
			interval += 5 * time.Second
		}

		time.Sleep(interval)
	}
}

func parseDeviceFlowError(resp *http.Response) (*DeviceFlowError, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read error response (HTTP %d): %w", resp.StatusCode, err)
	}
	var deviceErr DeviceFlowError
	if err := json.Unmarshal(body, &deviceErr); err != nil {
		return nil, fmt.Errorf("device authorization failed (HTTP %d): %s", resp.StatusCode, string(body))
	}
	return &deviceErr, nil
}

func RunOIDCLogin(opts OIDCFlowOptions) (*TokenResponse, error) {
	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	if opts.Scopes == "" {
		opts.Scopes = defaultOIDCScopes
	}

	discovery, err := oidcDiscoveryFunc(opts.IssuerURL, client)
	if err != nil {
		return nil, err
	}

	if discovery.DeviceAuthorizationEndpoint == "" {
		return nil, fmt.Errorf("OIDC provider at %s does not support device authorization", opts.IssuerURL)
	}

	deviceResp, err := deviceAuthorizationFunc(discovery.DeviceAuthorizationEndpoint, opts.ClientID, opts.Scopes, client)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nTo complete login, open the following URL in your browser:\n\n")
	fmt.Printf("  %s\n\n", deviceResp.VerificationURI)
	fmt.Printf("And enter the code: %s\n\n", deviceResp.UserCode)
	fmt.Printf("Waiting for authorization...\n")

	pollInterval := time.Duration(deviceResp.Interval) * time.Second
	if pollInterval == 0 {
		pollInterval = defaultPollInterval
	}

	tokenResp, err := pollTokenFunc(discovery.TokenEndpoint, deviceResp.DeviceCode, opts.ClientID, pollInterval, client)
	if err != nil {
		return nil, err
	}

	return tokenResp, nil
}
