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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	helpers "github.com/goharbor/harbor-cli/test/helper"
	"github.com/stretchr/testify/assert"
)

func newMockHarborServer(t *testing.T, validCreds map[string]string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || validCreds[username] != password {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"errors":[{"code":"UNAUTHORIZED","message":"unauthorized"}]}`)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"user_id":1,"username":"`+username+`"}`)
	}))
}

func Test_Login_Success(t *testing.T) {
	helpers.SetMockKeyring(t)
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)
	cmd := root.LoginCommand()

	server := newMockHarborServer(t, map[string]string{
		"harbor-cli": "Harbor12345",
	})
	defer server.Close()

	validServerAddresses := []string{
		server.URL,
		server.URL + "/",
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
	helpers.SetMockKeyring(t)
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
	helpers.SetMockKeyring(t)
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	server := newMockHarborServer(t, map[string]string{
		"harbor-cli": "Harbor12345",
	})
	defer server.Close()

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{server.URL})

	assert.NoError(t, cmd.Flags().Set("username", "does-not-exist"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for wrong username")
}

func Test_Login_Failure_WrongPassword(t *testing.T) {
	helpers.SetMockKeyring(t)
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	server := newMockHarborServer(t, map[string]string{
		"admin": "Harbor12345",
	})
	defer server.Close()

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{server.URL})

	assert.NoError(t, cmd.Flags().Set("username", "admin"))
	assert.NoError(t, cmd.Flags().Set("password", "wrong"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for wrong password")
}

func Test_Login_Success_RobotAccount(t *testing.T) {
	helpers.SetMockKeyring(t)
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	server := newMockHarborServer(t, map[string]string{
		"robot_harbor-cli": "Harbor12345",
	})
	defer server.Close()

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{server.URL})

	assert.NoError(t, cmd.Flags().Set("username", "robot_harbor-cli"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

	err := cmd.Execute()
	assert.NoError(t, err, "Expected no error for robot account login")
}
