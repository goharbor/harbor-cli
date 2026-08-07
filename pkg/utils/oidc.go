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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	oidcCLILoginPath   = "/c/oidc/login"
	oidcCLITokenPath   = "/c/oidc/cli-token" //nolint:gosec // endpoint is user-provided Harbor server URL for login.
	oidcCLIRefreshPath = "/c/oidc/refresh"
)

type OIDCLoginResponse struct {
	RedirectURL string `json:"redirect_url"`
	PollToken   string `json:"poll_token"`
}

type OIDCPollResponse struct {
	Status       string `json:"status"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Username     string `json:"username,omitempty"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
	Error        string `json:"error,omitempty"`
}

type OIDCRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type OIDCRefreshResponse struct {
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresAt    int64  `json:"expires_at"`
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
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OIDC login endpoint: %w", err)
	}
	q := u.Query()
	q.Set("mode", "cli")
	u.RawQuery = q.Encode()
	log.Debugf("initiating Harbor CLI OIDC login against %s", u.String())

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(u.String()) //nolint:gosec // endpoint is user-provided Harbor server URL for login.
	if err != nil {
		return nil, fmt.Errorf("failed to initiate OIDC login: %w", err)
	}
	defer resp.Body.Close()
	log.Debugf("received OIDC login response status %d from %s", resp.StatusCode, u.String())

	if resp.StatusCode >= http.StatusMultipleChoices && resp.StatusCode < http.StatusBadRequest {
		location := resp.Header.Get("Location")
		if location != "" {
			log.Debugf("Harbor returned browser redirect for CLI OIDC login: %s", location)
			return nil, fmt.Errorf("This Harbor instance may not support CLI OIDC login yet")
		}
		log.Debug("Harbor returned an HTTP redirect instead of a CLI OIDC JSON response")
		return nil, fmt.Errorf("Harbor did not return a CLI OIDC login response and issued an HTTP redirect instead. This Harbor instance may not support CLI OIDC login yet")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("failed to initiate OIDC login: status %d: %s", resp.StatusCode, string(body))
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "" && !strings.Contains(strings.ToLower(contentType), "application/json") {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		log.Debugf("Harbor returned non-JSON OIDC login response with content-type %q", contentType)
		return nil, fmt.Errorf("Harbor did not return a CLI OIDC login JSON response (content-type %q). This Harbor instance may not support CLI OIDC login yet: %s", contentType, string(body))
	}

	var loginResp OIDCLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("failed to decode OIDC login response: %w. This Harbor instance may not support CLI OIDC login yet", err)
	}
	if loginResp.RedirectURL == "" || loginResp.PollToken == "" {
		return nil, fmt.Errorf("invalid OIDC login response: missing redirect_url or poll_token")
	}
	log.Debug("received Harbor CLI OIDC login payload successfully")
	return &loginResp, nil
}

func PollForOIDCToken(serverAddress, pollToken string, timeout time.Duration) (*OIDCPollResponse, error) {
	if pollToken == "" {
		return nil, fmt.Errorf("poll token is required")
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
	q.Set("poll_token", pollToken)
	u.RawQuery = q.Encode()
	log.Debugf("starting Harbor CLI OIDC polling against %s with timeout %s", u.String(), timeout)

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		result, ready, err := pollOIDCTokenOnce(u.String())
		if err != nil {
			return nil, err
		}
		if ready {
			log.Debug("Harbor CLI OIDC polling completed with ready token response")
			return result, nil
		}
		if time.Now().After(deadline) {
			log.Debug("Harbor CLI OIDC polling timed out while waiting for authentication")
			return nil, fmt.Errorf("timed out waiting for OIDC authentication")
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, fmt.Errorf("timed out waiting for OIDC authentication")
		}
		select {
		case <-ticker.C:
			log.Debug("Harbor CLI OIDC token still pending; retrying poll")
		case <-time.After(remaining):
			log.Debug("Harbor CLI OIDC polling deadline reached while waiting for next retry")
			return nil, fmt.Errorf("timed out waiting for OIDC authentication")
		}
	}
}

func pollOIDCTokenOnce(endpoint string) (*OIDCPollResponse, bool, error) {
	resp, err := (&http.Client{Timeout: 30 * time.Second}).Get(endpoint) //nolint:gosec // endpoint is the Harbor server URL validated by PollForOIDCToken.
	if err != nil {
		return nil, false, fmt.Errorf("failed to poll OIDC token: %w", err)
	}
	defer resp.Body.Close()

	var pollResp OIDCPollResponse
	switch resp.StatusCode {
	case http.StatusAccepted:
		log.Debug("Harbor CLI OIDC poll returned pending status")
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
		log.Debugf("Harbor CLI OIDC poll returned ready status for user %s", pollResp.Username)
		return &pollResp, true, nil
	case http.StatusBadRequest:
		if err := json.NewDecoder(resp.Body).Decode(&pollResp); err != nil {
			return nil, false, fmt.Errorf("OIDC authentication failed")
		}
		if pollResp.Error != "" {
			log.Debugf("Harbor CLI OIDC poll returned failed status: %s", pollResp.Error)
			return nil, false, fmt.Errorf("OIDC authentication failed: %s", pollResp.Error)
		}
		log.Debug("Harbor CLI OIDC poll returned failed status without detailed error")
		return nil, false, fmt.Errorf("OIDC authentication failed")
	case http.StatusGone:
		if err := json.NewDecoder(resp.Body).Decode(&pollResp); err != nil {
			return nil, false, fmt.Errorf("OIDC login expired before token retrieval")
		}
		log.Debug("Harbor CLI OIDC poll returned expired status")
		return nil, false, fmt.Errorf("OIDC login expired before token retrieval")
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		log.Debugf("Harbor CLI OIDC poll returned unexpected status %d", resp.StatusCode)
		return nil, false, fmt.Errorf("failed to poll OIDC token: status %d: %s", resp.StatusCode, string(body))
	}
}

func RefreshOIDCToken(serverAddress, refreshToken string) (*OIDCRefreshResponse, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	serverAddress = FormatUrl(serverAddress)
	if err := ValidateURL(serverAddress); err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}

	endpoint, err := joinServerPath(serverAddress, oidcCLIRefreshPath)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(&OIDCRefreshRequest{RefreshToken: refreshToken})
	if err != nil {
		return nil, fmt.Errorf("failed to encode OIDC refresh request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh OIDC token: %w", err)
	}
	defer resp.Body.Close()

	var refreshResp OIDCRefreshResponse
	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
			return nil, fmt.Errorf("failed to decode OIDC refresh response: %w", err)
		}
		if refreshResp.IDToken == "" {
			return nil, fmt.Errorf("invalid OIDC refresh response: missing id_token")
		}
		return &refreshResp, nil
	case http.StatusBadRequest:
		if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err == nil && refreshResp.Error != "" {
			return nil, fmt.Errorf("OIDC refresh failed: %s", refreshResp.Error)
		}
		fallthrough
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("failed to refresh OIDC token: status %d: %s", resp.StatusCode, string(body))
	}
}

func joinServerPath(serverAddress, path string) (string, error) {
	u, err := url.Parse(serverAddress)
	if err != nil {
		return "", fmt.Errorf("failed to parse server URL: %w", err)
	}
	basePath := u.Path
	for len(basePath) > 0 && basePath[len(basePath)-1] == '/' {
		basePath = basePath[:len(basePath)-1]
	}
	u.Path = basePath + path
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}
