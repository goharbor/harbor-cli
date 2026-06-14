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
	"time"
)

const (
	oidcCLILoginPath = "/c/oidc/cli/login"
	oidcCLITokenPath = "/c/oidc/cli-token"
)

type OIDCLoginResponse struct {
	RedirectURL    string `json:"redirect_url"`
	TransactionID string `json:"transaction_id"`
}

type OIDCPollResponse struct {
	Status       string `json:"status"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Username     string `json:"username,omitempty"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
	Error        string `json:"error,omitempty"`
}

func InitiateOIDCLogin(serverAddress string) (*OIDCLoginResponse, error) {
	serverAddress = FormatUrl(serverAddress)
	if err := ValidateURL(serverAddress); err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}

	endpoint, err := joinServerPath(serverAddress, oidcCLILoginPath)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(endpoint) //nolint:gosec // endpoint is user-provided Harbor server URL for login.
	if err != nil {
		return nil, fmt.Errorf("failed to initiate OIDC login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("failed to initiate OIDC login: status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp OIDCLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("failed to decode OIDC login response: %w", err)
	}
	if loginResp.RedirectURL == "" || loginResp.TransactionID == "" {
		return nil, fmt.Errorf("invalid OIDC login response: missing redirect_url or transaction_id")
	}
	return &loginResp, nil
}

func PollForOIDCToken(serverAddress, transactionID string, timeout time.Duration) (*OIDCPollResponse, error) {
	if transactionID == "" {
		return nil, fmt.Errorf("transaction ID is required")
	}
	serverAddress = FormatUrl(serverAddress)
	if err := ValidateURL(serverAddress); err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}

	endpoint, err := joinServerPath(serverAddress, oidcCLITokenPath)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OIDC token endpoint: %w", err)
	}
	q := u.Query()
	q.Set("transaction_id", transactionID)
	u.RawQuery = q.Encode()

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		result, ready, err := pollOIDCTokenOnce(u.String())
		if err != nil {
			return nil, err
		}
		if ready {
			return result, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out waiting for OIDC authentication")
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, fmt.Errorf("timed out waiting for OIDC authentication")
		}
		select {
		case <-ticker.C:
		case <-time.After(remaining):
			return nil, fmt.Errorf("timed out waiting for OIDC authentication")
		}
	}
}

func pollOIDCTokenOnce(endpoint string) (*OIDCPollResponse, bool, error) {
	resp, err := http.Get(endpoint) //nolint:gosec // endpoint is the Harbor server URL validated by PollForOIDCToken.
	if err != nil {
		return nil, false, fmt.Errorf("failed to poll OIDC token: %w", err)
	}
	defer resp.Body.Close()

	var pollResp OIDCPollResponse
	switch resp.StatusCode {
	case http.StatusAccepted:
		return &OIDCPollResponse{Status: "pending"}, false, nil
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&pollResp); err != nil {
			return nil, false, fmt.Errorf("failed to decode OIDC token response: %w", err)
		}
		if pollResp.Status != "ready" {
			return nil, false, fmt.Errorf("unexpected OIDC token status: %s", pollResp.Status)
		}
		if pollResp.IDToken == "" || pollResp.Username == "" {
			return nil, false, fmt.Errorf("invalid OIDC token response: missing id_token or username")
		}
		return &pollResp, true, nil
	case http.StatusBadRequest:
		if err := json.NewDecoder(resp.Body).Decode(&pollResp); err != nil {
			return nil, false, fmt.Errorf("OIDC authentication failed")
		}
		if pollResp.Error != "" {
			return nil, false, fmt.Errorf("OIDC authentication failed: %s", pollResp.Error)
		}
		return nil, false, fmt.Errorf("OIDC authentication failed")
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, false, fmt.Errorf("failed to poll OIDC token: status %d: %s", resp.StatusCode, string(body))
	}
}

func joinServerPath(serverAddress, path string) (string, error) {
	u, err := url.Parse(serverAddress)
	if err != nil {
		return "", fmt.Errorf("failed to parse server URL: %w", err)
	}
	u.Path = path
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}
