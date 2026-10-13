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
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/goharbor/harbor-cli/pkg/utils"
	helpers "github.com/goharbor/harbor-cli/test/helper"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Logins against a real Harbor live in test/e2e; these cover only what can be
// decided without a server.

// unreachable is a port nothing listens on, so dialing it fails immediately
// without leaving the machine.
const unreachable = "http://127.0.0.1:1"

func Test_Login_Failure_UnreachableServer(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{unreachable})

	assert.NoError(t, cmd.Flags().Set("username", "harbor-cli"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for unreachable server")
}

func Test_Login_ExistingCredentials(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)
	helpers.SetMockKeyring(t)

	var logs bytes.Buffer
	log.SetOutput(&logs)
	defer log.SetOutput(os.Stderr)

	// Each login runs in a fresh process in real use, so reload the config
	// written by the previous login.
	reloadConfig := func() {
		utils.ConfigInitialization = &utils.Once{}
		utils.CurrentHarborConfig = nil
		utils.InitConfig(data.ConfigPath, true)
	}
	login := func(username, password string) {
		reloadConfig()
		cmd := root.LoginCommand()
		cmd.SetArgs([]string{unreachable})
		assert.NoError(t, cmd.Flags().Set("username", username))
		assert.NoError(t, cmd.Flags().Set("password", password))
		assert.NoError(t, cmd.Flags().Set("skip-verify-client", "true"))
		assert.NoError(t, cmd.Execute())
	}

	login("harbor-cli", "Harbor12345")
	login("other-user", "Harbor12345")

	// The same credentials are not stored again, but their context becomes current
	logs.Reset()
	login("harbor-cli", "Harbor12345")
	assert.Contains(t, logs.String(), "They were not added again")
	reloadConfig()
	config, err := utils.GetCurrentHarborConfig()
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(config.CurrentCredentialName, "harbor-cli@"), config.CurrentCredentialName)

	logs.Reset()
	login("harbor-cli", "Harbor67890")
	assert.Contains(t, logs.String(), "the password is different")
}

// Keep this last: --password-stdin is bound to a package level variable that
// stays set for the rest of the process once a test flips it.
func Test_Login_Failure_MutuallyExclusiveFlags(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	cmd := root.LoginCommand()
	cmd.SetArgs([]string{unreachable})

	assert.NoError(t, cmd.Flags().Set("username", "admin"))
	assert.NoError(t, cmd.Flags().Set("password", "Harbor12345"))
	assert.NoError(t, cmd.Flags().Set("password-stdin", "true"))

	err := cmd.Execute()
	assert.Error(t, err, "Expected error when both --password and --password-stdin are set")
}
