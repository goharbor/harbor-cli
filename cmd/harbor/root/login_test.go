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
package root_test

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/goharbor/harbor-cli/pkg/utils"
	helpers "github.com/goharbor/harbor-cli/test/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testEncryptionKey = "12345678901234567890123456789012"

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

func newFakeHarborTransport(t *testing.T, authMode string) http.RoundTripper {
	t.Helper()

	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v2.0/systeminfo":
			return jsonResponse(t, r, http.StatusOK, map[string]string{
				"auth_mode": authMode,
			}), nil
		case "/api/v2.0/users/current":
			return fakeUserInfoResponse(t, r), nil
		case "/api/v2.0/projects":
			return fakeProjectsResponse(t, r), nil
		case "/api/v2.0/ping":
			return fakePingResponse(r), nil
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("not found")),
				Request:    r,
			}, nil
		}
	})
}

func fakeUserInfoResponse(t *testing.T, r *http.Request) *http.Response {
	t.Helper()

	username, password, ok := basicAuth(r)
	if !ok {
		return harborErrorResponse(r, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	}

	switch {
	case username == "harbor-cli" && password == "Harbor12345":
		return jsonResponse(t, r, http.StatusOK, map[string]any{"username": username})
	case username == "robot_harbor-cli" && password == "Harbor12345":
		return harborErrorResponse(r, http.StatusPreconditionFailed, "PRECONDITION_FAILED", "precondition failed")
	default:
		return harborErrorResponse(r, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	}
}

func fakeProjectsResponse(t *testing.T, r *http.Request) *http.Response {
	t.Helper()

	username, password, ok := basicAuth(r)
	if !ok {
		return harborErrorResponse(r, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	}

	switch {
	case username == "harbor-cli" && password == "Harbor12345":
		return jsonResponse(t, r, http.StatusOK, []map[string]any{{"name": "library"}})
	case username == "robot_harbor-cli" && password == "Harbor12345":
		return jsonResponse(t, r, http.StatusOK, []map[string]any{{"name": "robot-project"}})
	default:
		return harborErrorResponse(r, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	}
}

func fakePingResponse(r *http.Request) *http.Response {
	username, password, ok := basicAuth(r)
	if !ok {
		return harborErrorResponse(r, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	}

	switch {
	case username == "harbor-cli" && password == "Harbor12345":
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("")), Request: r}
	case username == "robot_harbor-cli" && password == "Harbor12345":
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("")), Request: r}
	default:
		return harborErrorResponse(r, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	}
}

func basicAuth(r *http.Request) (string, string, bool) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", false
	}

	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(raw), ":")
	if !ok {
		return "", "", false
	}
	return username, password, true
}

func harborErrorResponse(r *http.Request, status int, code, message string) *http.Response {
	payload, _ := json.Marshal(map[string]any{
		"errors": []map[string]string{{
			"code":    code,
			"message": message,
		}},
	})

	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(string(payload))),
		Request:    r,
	}
}

func jsonResponse(t *testing.T, r *http.Request, status int, payload any) *http.Response {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(string(body))),
		Request:    r,
	}
}

func Test_Login_Success(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SafeSetEnv("HARBOR_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString([]byte(testEncryptionKey)))
	t.Cleanup(func() { helpers.SafeUnsetEnv("HARBOR_ENCRYPTION_KEY") })
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"https://harbor.example.com"})

	assert.NoError(t, cmd.Flags().Set("username", "harbor-cli"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))
	assert.NoError(t, cmd.Flags().Set("skip-verify-client", "true"))

	err := cmd.Execute()
	assert.NoError(t, err)
}

func Test_Login_Failure_WrongServer(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"wrongserver"})

	assert.NoError(t, cmd.Flags().Set("username", "harbor-cli"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for invalid server")
}

func Test_Login_StoresCredentialsWhenVerificationSkipped(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SafeSetEnv("HARBOR_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString([]byte(testEncryptionKey)))
	t.Cleanup(func() { helpers.SafeUnsetEnv("HARBOR_ENCRYPTION_KEY") })
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"https://harbor.example.com"})

	assert.NoError(t, cmd.Flags().Set("username", "does-not-exist"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))
	assert.NoError(t, cmd.Flags().Set("skip-verify-client", "true"))

	err := cmd.Execute()
	assert.NoError(t, err)

	cred, err := utils.GetCredentials(utils.DefaultCredentialName("does-not-exist", "https://harbor.example.com"))
	assert.NoError(t, err)
	assert.Equal(t, "does-not-exist", cred.Username)
}

func Test_Login_Success_RobotAccount(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SafeSetEnv("HARBOR_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString([]byte(testEncryptionKey)))
	t.Cleanup(func() { helpers.SafeUnsetEnv("HARBOR_ENCRYPTION_KEY") })
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"https://harbor.example.com"})

	assert.NoError(t, cmd.Flags().Set("username", "robot_harbor-cli"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))
	assert.NoError(t, cmd.Flags().Set("skip-verify-client", "true"))

	err := cmd.Execute()
	assert.NoError(t, err, "Expected no error for robot account login")
}

func Test_Login_Failure_MutuallyExclusiveFlags(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"http://demo.goharbor.io"})

	assert.NoError(t, cmd.Flags().Set("username", "admin"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))
	assert.NoError(t, cmd.Flags().Set("password-stdin", "true"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error when both --password and --password-stdin are set")
}

func Test_Login_Failure_InvalidAuthModeForOIDCHarbor(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)
	withDefaultTransport(t, newFakeHarborTransport(t, "oidc_auth"))

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"https://harbor.example.com"})

	assert.NoError(t, cmd.Flags().Set("auth-mode", "ldap"))
	assert.NoError(t, cmd.Flags().Set("username", "admin"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

	err := cmd.Execute()
	assert.ErrorContains(t, err, "LDAP login is not available because Harbor auth_mode is oidc_auth")
}

func Test_RunOIDCLogin_Failure_MissingServer(t *testing.T) {
	err := root.RunOIDCLogin("")
	assert.Error(t, err)
}

func Test_Login_AutoDetectsOIDCAuthMode(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SafeSetEnv("HARBOR_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString([]byte(testEncryptionKey)))
	t.Cleanup(func() { helpers.SafeUnsetEnv("HARBOR_ENCRYPTION_KEY") })
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v2.0/systeminfo":
			return jsonResponse(t, r, http.StatusOK, map[string]string{
				"auth_mode": "oidc_auth",
			}), nil
		case "/c/oidc/login":
			assert.Equal(t, "cli", r.URL.Query().Get("mode"))
			return jsonResponse(t, r, http.StatusOK, utils.OIDCLoginResponse{
				RedirectURL: "https://idp.example/authorize",
				PollToken:   "poll-token-1",
			}), nil
		case "/c/oidc/cli-token":
			assert.Equal(t, "poll-token-1", r.URL.Query().Get("poll_token"))
			return jsonResponse(t, r, http.StatusOK, utils.OIDCPollResponse{
				Status:   "ready",
				IDToken:  "id-token",
				Username: "alice",
			}), nil
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("not found")),
				Request:    r,
			}, nil
		}
	}))

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"https://harbor.example.com"})

	err := cmd.Execute()
	assert.NoError(t, err)

	cred, err := utils.GetCredentials(utils.DefaultCredentialName("alice", "https://harbor.example.com"))
	assert.NoError(t, err)
	assert.Equal(t, utils.AuthTypeOIDC, cred.AuthType)
}

func Test_Login_AllowsDBAuthModeWhenHarborUsesOIDC(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SafeSetEnv("HARBOR_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString([]byte(testEncryptionKey)))
	t.Cleanup(func() { helpers.SafeUnsetEnv("HARBOR_ENCRYPTION_KEY") })
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)
	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/api/v2.0/systeminfo" {
			return jsonResponse(t, r, http.StatusOK, map[string]string{"auth_mode": "oidc_auth"}), nil
		}
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("not found")),
			Request:    r,
		}, nil
	}))

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"https://harbor.example.com"})

	assert.NoError(t, cmd.Flags().Set("auth-mode", "db"))
	assert.NoError(t, cmd.Flags().Set("username", "alice"))
	assert.NoError(t, cmd.Flags().Set("password", "cli-secret"))
	assert.NoError(t, cmd.Flags().Set("skip-verify-client", "true"))

	err := cmd.Execute()
	assert.NoError(t, err)

	cred, err := utils.GetCredentials(utils.DefaultCredentialName("alice", "https://harbor.example.com"))
	assert.NoError(t, err)
	assert.Equal(t, "alice", cred.Username)
	assert.Empty(t, cred.AuthType)
}

func Test_RunOIDCLogin_Success(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SafeSetEnv("HARBOR_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString([]byte(testEncryptionKey)))
	t.Cleanup(func() { helpers.SafeUnsetEnv("HARBOR_ENCRYPTION_KEY") })
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	withDefaultTransport(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/c/oidc/login":
			assert.Equal(t, "cli", r.URL.Query().Get("mode"))
			return jsonResponse(t, r, http.StatusOK, utils.OIDCLoginResponse{
				RedirectURL: "https://idp.example/authorize",
				PollToken:   "poll-token-1",
			}), nil
		case "/c/oidc/cli-token":
			assert.Equal(t, "poll-token-1", r.URL.Query().Get("poll_token"))
			return jsonResponse(t, r, http.StatusOK, utils.OIDCPollResponse{
				Status:   "ready",
				IDToken:  "id-token",
				Username: "alice",
			}), nil
		default:
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("not found")),
				Request:    r,
			}, nil
		}
	}))

	err := root.RunOIDCLogin("https://harbor.example.com")
	assert.NoError(t, err)

	cred, err := utils.GetCredentials(utils.DefaultCredentialName("alice", "https://harbor.example.com"))
	assert.NoError(t, err)
	assert.Equal(t, utils.AuthTypeOIDC, cred.AuthType)
}
