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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/harbor"
	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	log "github.com/sirupsen/logrus"
)

var (
	ClientInstance *v2client.HarborAPI
	ClientOnce     sync.Once
	ClientErr      error
)

const oidcRefreshFailureMessage = "Unable to refresh OIDC session. Please try again or run 'harbor login <server> --oidc'."

const oidcRetryHeader = "X-Harbor-CLI-OIDC-Retry"

type oidcTokenManager struct {
	mu         sync.RWMutex
	credential Credential
	token      string
	refreshFn  func(Credential) (string, error)
}

type oidcRetryTransport struct {
	base         http.RoundTripper
	tokenManager *oidcTokenManager
}

func GetClient() (*v2client.HarborAPI, error) {
	ClientOnce.Do(func() {
		config, err := GetCurrentHarborConfig()
		if err != nil {
			ClientErr = fmt.Errorf("failed to get current credential name: %v", err)
			return
		}
		credentialName := config.CurrentCredentialName
		if credentialName == "" {
			ClientErr = fmt.Errorf("no Harbor credentials found. Please run `harbor login` to configure access")
			return
		}

		ClientInstance, ClientErr = GetClientByCredentialName(credentialName)
		if ClientErr != nil {
			log.Errorf("failed to initialize client: %v", ClientErr)
			return
		}
	})

	return ClientInstance, ClientErr
}

func ContextWithClient() (context.Context, *v2client.HarborAPI, error) {
	client, err := GetClient()
	if err != nil {
		return nil, nil, err
	}
	ctx := context.Background()
	return ctx, client, nil
}

func GetClientByConfig(clientConfig *harbor.ClientSetConfig) *v2client.HarborAPI {
	cs, err := harbor.NewClientSet(clientConfig)
	if err != nil {
		panic(err)
	}
	return cs.V2()
}

// Returns Harbor v2 client after resolving the credential name
func GetClientByCredentialName(credentialName string) (*v2client.HarborAPI, error) {
	credential, err := GetCredentials(credentialName)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential %s: %w", credentialName, err)
	}
	if credential.AuthType == AuthTypeOIDC {
		return getOIDCClient(credential)
	}

	// Get encryption key
	key, err := GetEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption key: %w", err)
	}

	// Decrypt password
	decryptedPassword, err := Decrypt(key, string(credential.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	clientConfig := &harbor.ClientSetConfig{
		URL:      credential.ServerAddress,
		Username: credential.Username,
		Password: decryptedPassword,
	}
	return GetClientByConfig(clientConfig), nil
}

func getOIDCClient(credential Credential) (*v2client.HarborAPI, error) {
	idToken, err := GetDecryptedIDToken(credential.Name)
	if err != nil {
		return nil, err
	}

	if oidcTokenNeedsRefresh(credential, idToken) {
		idToken, err = refreshOIDCCredential(credential)
		if err != nil {
			return nil, err
		}
	}

	return buildOIDCClient(credential, idToken)
}

func refreshOIDCCredential(credential Credential) (string, error) {
	refreshToken, err := GetDecryptedRefreshToken(credential.Name)
	if err != nil {
		return "", fmt.Errorf("failed to load OIDC refresh token: %w", err)
	}
	if refreshToken == "" {
		return "", fmt.Errorf(oidcRefreshFailureMessage)
	}

	refreshResp, err := RefreshOIDCToken(credential.ServerAddress, refreshToken)
	if err != nil {
		log.WithError(err).Warn("failed to refresh OIDC token")
		return "", fmt.Errorf(oidcRefreshFailureMessage)
	}

	nextRefreshToken := refreshResp.RefreshToken
	if nextRefreshToken == "" {
		nextRefreshToken = refreshToken
	}

	harborData, err := GetCurrentHarborData()
	if err != nil {
		return "", fmt.Errorf("failed to get current Harbor data: %w", err)
	}
	if err := UpdateOIDCTokens(credential.Name, refreshResp.IDToken, nextRefreshToken, refreshResp.ExpiresAt, harborData.ConfigPath); err != nil {
		return "", fmt.Errorf("failed to persist refreshed OIDC tokens: %w", err)
	}

	return refreshResp.IDToken, nil
}

func oidcTokenNeedsRefresh(credential Credential, idToken string) bool {
	expiresAt, err := oidcTokenExpiryUnix(idToken)
	if err != nil {
		log.Debugf("failed to parse OIDC token expiry from JWT, falling back to stored expires_at: %v", err)
		expiresAt = credential.ExpiresAt
	}
	if expiresAt <= 0 {
		return false
	}
	return time.Now().Unix() >= expiresAt-60
}

func oidcTokenExpiryUnix(idToken string) (int64, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid JWT format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, fmt.Errorf("failed to unmarshal JWT payload: %w", err)
	}
	if claims.Exp <= 0 {
		return 0, fmt.Errorf("JWT exp claim is missing")
	}
	return claims.Exp, nil
}

func buildOIDCClient(credential Credential, idToken string) (*v2client.HarborAPI, error) {
	tokenManager := &oidcTokenManager{
		credential: credential,
		token:      idToken,
		refreshFn:  refreshOIDCCredential,
	}

	return buildClientWithAuth(credential.ServerAddress, tokenManager, &oidcRetryTransport{
		base:         http.DefaultTransport,
		tokenManager: tokenManager,
	})
}

func buildClientWithAuth(serverAddress string, tokenManager *oidcTokenManager, transport http.RoundTripper) (*v2client.HarborAPI, error) {
	u, err := url.Parse(serverAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server URL: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid server URL: %s", serverAddress)
	}

	cfg := &harbor.Config{
		URL:       u,
		Transport: transport,
		AuthInfo: runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) error {
			return req.SetHeaderParam("Authorization", "Bearer "+tokenManager.Token())
		}),
	}

	return v2client.New(cfg.ToV2Config()), nil
}

func (m *oidcTokenManager) Token() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.token
}

func (m *oidcTokenManager) Refresh() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	refreshFn := m.refreshFn
	if refreshFn == nil {
		refreshFn = refreshOIDCCredential
	}

	token, err := refreshFn(m.credential)
	if err != nil {
		return "", err
	}
	m.token = token
	return token, nil
}

func (t *oidcRetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	resp, err := base.RoundTrip(req)
	if err != nil || resp == nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}
	if req.Header.Get(oidcRetryHeader) == "1" || !canRetryOIDCRequest(req) {
		return resp, nil
	}

	if _, err := t.tokenManager.Refresh(); err != nil {
		drainAndCloseResponse(resp)
		return nil, err
	}

	drainAndCloseResponse(resp)

	retryReq, err := cloneRequestForRetry(req)
	if err != nil {
		return nil, err
	}
	retryReq.Header.Set("Authorization", "Bearer "+t.tokenManager.Token())
	retryReq.Header.Set(oidcRetryHeader, "1")

	return base.RoundTrip(retryReq)
}

func canRetryOIDCRequest(req *http.Request) bool {
	if req == nil {
		return false
	}
	if req.Body == nil || req.Body == http.NoBody {
		return true
	}
	return req.GetBody != nil
}

func cloneRequestForRetry(req *http.Request) (*http.Request, error) {
	retryReq := req.Clone(req.Context())
	if req.Body != nil && req.Body != http.NoBody {
		if req.GetBody == nil {
			return nil, fmt.Errorf("request body cannot be retried")
		}
		body, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("failed to reset request body: %w", err)
		}
		retryReq.Body = body
	} else {
		retryReq.Body = http.NoBody
	}
	retryReq.Header = req.Header.Clone()
	return retryReq, nil
}

func drainAndCloseResponse(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	_ = resp.Body.Close()
}
