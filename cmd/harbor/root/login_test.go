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
	"os"
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
	harborURL := "https://demo.goharbor.io"
	if url := os.Getenv("HARBOR_URL"); url != "" {
		harborURL = url
	}

	validServerAddresses := []string{
		harborURL,
	}

	for _, serverAddress := range validServerAddresses {
		t.Run("ValidServer_"+serverAddress, func(t *testing.T) {
			args := []string{serverAddress}
			cmd.SetArgs(args)

			username := "harbor-cli"
			if u := os.Getenv("HARBOR_USERNAME"); u != "" {
				username = u
			}
			password := "Harbor12345"
			if p := os.Getenv("HARBOR_PASSWORD"); p != "" {
				password = p
			}

			assert.NoError(t, cmd.Flags().Set("username", username))
			assert.NoError(t, cmd.Flags().Set("password", password))

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
	harborURL := "http://demo.goharbor.io"
	if url := os.Getenv("HARBOR_URL"); url != "" {
		harborURL = url
	}
	cmd.SetArgs([]string{harborURL})

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
	harborURL := "http://demo.goharbor.io"
	if url := os.Getenv("HARBOR_URL"); url != "" {
		harborURL = url
	}
	cmd.SetArgs([]string{harborURL})

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
	harborURL := "https://demo.goharbor.io"
	if url := os.Getenv("HARBOR_URL"); url != "" {
		t.Skip("Skipping robot account test for custom HARBOR_URL")
	}
	cmd.SetArgs([]string{harborURL})

	assert.NoError(t, cmd.Flags().Set("username", "robot_harbor-cli"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

	err := cmd.Execute()
	assert.NoError(t, err, "Expected no error for robot account login")
}
