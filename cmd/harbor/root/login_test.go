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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	helpers "github.com/goharbor/harbor-cli/test/helper"
	"github.com/stretchr/testify/assert"
)

func Test_Login_Success(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)
	cmd := root.LoginCommand()
	validServerAddresses := []string{
		"http://demo.goharbor.io:80",
		"https://demo.goharbor.io:443",
		"http://demo.goharbor.io",
		"https://demo.goharbor.io",
	}

	for _, serverAddress := range validServerAddresses {
		t.Run("ValidServer_"+serverAddress, func(t *testing.T) {
			args := []string{serverAddress}
			cmd.SetArgs(args)

			assert.NoError(t, cmd.Flags().Set("username", "harbor-cli"))
			assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

			err := cmd.Execute()
			assert.NoError(t, err, "Expected no error for server: %s", serverAddress)
		})
	}
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

func Test_Login_Failure_WrongUsername(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"http://demo.goharbor.io"})

	assert.NoError(t, cmd.Flags().Set("username", "does-not-exist"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for wrong username")
}

func Test_Login_Failure_WrongPassword(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"http://demo.goharbor.io"})

	assert.NoError(t, cmd.Flags().Set("username", "admin"))
	assert.NoError(t, cmd.Flags().Set("password", "wrong"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for wrong password")
}

func Test_Login_Success_RobotAccount(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{"https://demo.goharbor.io"})

	assert.NoError(t, cmd.Flags().Set("username", "robot_harbor-cli"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

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

func Test_Login_Validation_MockServer(t *testing.T) {
	tests := []struct {
		name                   string
		currentUserStatus      int
		currentUserResponse    string
		projectsStatus         int
		projectsResponse       string
		pingStatus             int
		pingResponse           string
		expectError            bool
		expectedErrSubstr      string
	}{
		{
			name:                "Human Login Success",
			currentUserStatus:   http.StatusOK,
			currentUserResponse: `{"username": "human-user", "sysadmin_flag": false}`,
			expectError:         false,
		},
		{
			name:                "Invalid Credentials (401)",
			currentUserStatus:   http.StatusUnauthorized,
			currentUserResponse: `{"errors": [{"code": "UNAUTHORIZED", "message": "unauthorized"}]}`,
			expectError:         true,
			expectedErrSubstr:   "authentication failed, check your credentials",
		},
		{
			name:                "Robot Login Success (Standard fallback)",
			currentUserStatus:   http.StatusForbidden,
			currentUserResponse: `{"errors": [{"code": "FORBIDDEN", "message": "forbidden"}]}`,
			projectsStatus:      http.StatusOK,
			projectsResponse:    `[]`,
			pingStatus:          http.StatusOK,
			pingResponse:        `"pong"`,
			expectError:         false,
		},
		{
			name:                "Robot Login Success (Restricted project access 403 fallback)",
			currentUserStatus:   http.StatusForbidden,
			currentUserResponse: `{"errors": [{"code": "FORBIDDEN", "message": "forbidden"}]}`,
			projectsStatus:      http.StatusForbidden,
			projectsResponse:    `{"errors": [{"code": "FORBIDDEN", "message": "forbidden"}]}`,
			pingStatus:          http.StatusOK,
			pingResponse:        `"pong"`,
			expectError:         false,
		},
		{
			name:                "Robot Login Failure (Bad credentials 401 on fallback)",
			currentUserStatus:   http.StatusForbidden,
			currentUserResponse: `{"errors": [{"code": "FORBIDDEN", "message": "forbidden"}]}`,
			projectsStatus:      http.StatusUnauthorized,
			projectsResponse:    `{"errors": [{"code": "UNAUTHORIZED", "message": "unauthorized"}]}`,
			pingStatus:          http.StatusOK,
			pingResponse:        `"pong"`,
			expectError:         true,
			expectedErrSubstr:   "authentication failed, check your credentials",
		},
		{
			name:                "Server Error (500)",
			currentUserStatus:   http.StatusInternalServerError,
			currentUserResponse: `Internal Server Error`,
			projectsStatus:      http.StatusInternalServerError,
			projectsResponse:    `Internal Server Error`,
			pingStatus:          http.StatusInternalServerError,
			pingResponse:        `Internal Server Error`,
			expectError:         true,
			expectedErrSubstr:   "server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			data := helpers.Initialize(t, tempDir)
			defer helpers.ConfigCleanup(t, data)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/api/v2.0/users/current":
					w.WriteHeader(tt.currentUserStatus)
					w.Write([]byte(tt.currentUserResponse))
				case "/api/v2.0/projects":
					w.WriteHeader(tt.projectsStatus)
					w.Write([]byte(tt.projectsResponse))
				case "/api/v2.0/ping":
					w.WriteHeader(tt.pingStatus)
					w.Write([]byte(tt.pingResponse))
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			cmd := root.LoginCommand()
			cmd.SetArgs([]string{server.URL})

			assert.NoError(t, cmd.Flags().Set("username", "test-user"))
			assert.NoError(t, cmd.Flags().Set("password", "test-password"))

			err := cmd.Execute()
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrSubstr != "" {
					assert.Contains(t, err.Error(), tt.expectedErrSubstr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
