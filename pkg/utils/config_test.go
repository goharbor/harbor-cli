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
package utils_test

import (
	"path/filepath"
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/goharbor/harbor-cli/pkg/utils"

	helpers "github.com/goharbor/harbor-cli/test/helper"
	"github.com/stretchr/testify/assert"
)

func Test_Config_EnvVar(t *testing.T) {
	utils.ConfigInitialization.Reset() // Reset sync.Once for the test
	helpers.SetMockKeyring(t)
	tempDir := t.TempDir()
	helpers.SafeSetEnv("HARBOR_CLI_CONFIG", filepath.Join(tempDir, "config.yaml"))
	helpers.SafeSetEnv("XDG_DATA_HOME", filepath.Join(tempDir, ".data"))
	utils.InitConfig("", false)
	cds := root.RootCmd()
	err := cds.Execute()
	assert.NoError(t, err, "Expected no error for Root command")
	assert.NoError(t, err, "Expected no error for Root command execution")

	currentData, err := utils.GetCurrentHarborData()
	assert.NoError(t, err, "Expected no error when fetching HarborData")
	defer helpers.ConfigCleanup(t, currentData)

	currentConfig, err := utils.GetCurrentHarborConfig()
	assert.NoError(t, err, "Expected no error when fetching HarborConfig")
	assert.NotNil(t, currentConfig, "Configuration should not be nil")
	assert.NotNil(t, currentConfig.CurrentCredentialName, "CurrentCredentialName should not be nil")
	assert.NotNil(t, currentConfig.Credentials, "Credentials should not be nil")
	assert.NotNil(t, currentData.ConfigPath, "ConfigPath should not be nil")
}

func Test_Config_Vanilla(t *testing.T) {
	utils.ConfigInitialization.Reset() // Reset sync.Once for the test
	helpers.SetMockKeyring(t)
	utils.InitConfig("", false)
	cds := root.RootCmd()
	err := cds.Execute()
	assert.NoError(t, err, "Expected no error for Root command")
	assert.NoError(t, err, "Expected no error for Root command execution")
	currentData, err := utils.GetCurrentHarborData()
	assert.NoError(t, err, "Expected no error when fetching HarborData")
	defer helpers.ConfigCleanup(t, currentData)

	currentConfig, err := utils.GetCurrentHarborConfig()
	assert.NoError(t, err, "Expected no error when fetching HarborConfig")
	assert.NotNil(t, currentConfig, "Configuration should not be nil")
	assert.NotNil(t, currentConfig.CurrentCredentialName, "CurrentCredentialName should not be nil")
	assert.NotNil(t, currentConfig.Credentials, "Credentials should not be nil")
	assert.NotNil(t, currentData.ConfigPath, "ConfigPath should not be nil")
}

func Test_Config_Xdg(t *testing.T) {
	utils.ConfigInitialization.Reset() // Reset sync.Once for the test
	helpers.SetMockKeyring(t)
	tempDir := t.TempDir()
	helpers.SafeSetEnv("HARBOR_CLI_CONFIG", filepath.Join(tempDir, "config.yaml"))
	helpers.SafeSetEnv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))
	helpers.SafeSetEnv("XDG_DATA_HOME", filepath.Join(tempDir, ".data"))
	utils.InitConfig("", false)
	cds := root.RootCmd()
	err := cds.Execute()
	assert.NoError(t, err, "Expected no error for Root command")
	assert.NoError(t, err, "Expected no error for Root command execution")

	currentData, err := utils.GetCurrentHarborData()
	assert.NoError(t, err, "Expected no error when fetching HarborData")
	defer helpers.ConfigCleanup(t, currentData)

	currentConfig, err := utils.GetCurrentHarborConfig()
	assert.NoError(t, err, "Expected no error when fetching HarborConfig")
	assert.NotNil(t, currentConfig, "Configuration should not be nil")
	assert.NotNil(t, currentConfig.CurrentCredentialName, "CurrentCredentialName should not be nil")
	assert.NotNil(t, currentConfig.Credentials, "Credentials should not be nil")
	assert.NotNil(t, currentData.ConfigPath, "ConfigPath should not be nil")
}

func Test_Config_Flag(t *testing.T) {
	tempDir := t.TempDir()
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	testConfigFile := filepath.Join(tempDir, "config.yaml")
	utils.InitConfig(testConfigFile, true)
	currentConfig, err := utils.GetCurrentHarborConfig()
	assert.NoError(t, err, "Expected no error when fetching HarborConfig")
	assert.NotNil(t, currentConfig, "Configuration should not be nil")
	assert.NotNil(t, currentConfig.CurrentCredentialName, "CurrentCredentialName should not be nil")
	assert.NotNil(t, currentConfig.Credentials, "Credentials should not be nil")
	assert.NotNil(t, data.ConfigPath, "ConfigPath should not be nil")
}

func Test_AddOIDCCredentials(t *testing.T) {
	tempDir := t.TempDir()
	helpers.SetMockKeyring(t)
	data := helpers.Initialize(t, tempDir)
	defer helpers.ConfigCleanup(t, data)

	err := utils.AddOIDCCredentials("https://demo.goharbor.io", "alice", "id-token", "refresh-token", 12345, data.ConfigPath)
	assert.NoError(t, err)

	cred, err := utils.GetCredentials("alice@https-demo-goharbor-io")
	assert.NoError(t, err)
	assert.Equal(t, utils.AuthTypeOIDC, cred.AuthType)
	assert.Equal(t, "alice", cred.Username)
	assert.Equal(t, "https://demo.goharbor.io", cred.ServerAddress)
	assert.Equal(t, int64(12345), cred.ExpiresAt)
	assert.NotEmpty(t, cred.IDToken)
	assert.NotEmpty(t, cred.RefreshToken)
	assert.Empty(t, cred.Password)

	idToken, err := utils.GetDecryptedIDToken(cred.Name)
	assert.NoError(t, err)
	assert.Equal(t, "id-token", idToken)
}
