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
	"bytes"
	"io"
	"os"
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

// captureStdoutStderr redirects os.Stdout/os.Stderr while fn runs and returns
// whatever each stream received.
func captureStdoutStderr(t *testing.T, fn func()) (stdout string, stderr string) {
	t.Helper()

	origOut, origErr := os.Stdout, os.Stderr
	rOut, wOut, err := os.Pipe()
	assert.NoError(t, err)
	rErr, wErr, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout, os.Stderr = wOut, wErr
	defer func() { os.Stdout, os.Stderr = origOut, origErr }()

	fn()

	assert.NoError(t, wOut.Close())
	assert.NoError(t, wErr.Close())

	var bufOut, bufErr bytes.Buffer
	_, _ = io.Copy(&bufOut, rOut)
	_, _ = io.Copy(&bufErr, rErr)
	return bufOut.String(), bufErr.String()
}

// Test_CreateFiles_DoNotPolluteStdout is a regression test for #1053: the
// data/config file auto-creation notices must go to stderr, not stdout.
// `harbor completion bash` captures stdout into the generated completion
// script, so a stray "Config file created at ..."/"Data file created at ..."
// line on stdout ends up sourced by the shell as a command ("command not
// found: Config"/"command not found: Data").
func Test_CreateFiles_DoNotPolluteStdout(t *testing.T) {
	tempDir := t.TempDir()
	dataPath := filepath.Join(tempDir, ".data", "data.yaml")
	configPath := filepath.Join(tempDir, ".config", "config.yaml")

	stdout, stderr := captureStdoutStderr(t, func() {
		assert.NoError(t, utils.CreateConfigFile(configPath))
		assert.NoError(t, utils.CreateDataFile(dataPath, configPath))
	})

	assert.NotContains(t, stdout, "Config file created")
	assert.NotContains(t, stdout, "Data file created")
	assert.Empty(t, stdout, "file-creation notices must not be written to stdout")

	// The notices are still surfaced to the user, just on stderr.
	assert.Contains(t, stderr, "Config file created at "+configPath)
	assert.Contains(t, stderr, "Data file created at "+dataPath)
}
