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
	"os"
	"path/filepath"
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"

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

func Test_Config_Flag_Directory_Rejected(t *testing.T) {
	configDir := filepath.Join(t.TempDir(), "harbor-cli")
	assert.NoError(t, os.MkdirAll(configDir, os.ModePerm))

	_, err := utils.DetermineConfigPath(configDir, true)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "points to a directory")
	}
}

func Test_Config_EnvVar_Directory_Rejected(t *testing.T) {
	configDir := filepath.Join(t.TempDir(), "harbor-cli")
	assert.NoError(t, os.MkdirAll(configDir, os.ModePerm))
	t.Setenv("HARBOR_CLI_CONFIG", configDir)

	_, err := utils.DetermineConfigPath("", false)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "points to a directory")
	}
}

func Test_Config_Flag_ExistingFile_Accepted(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yaml")
	assert.NoError(t, os.WriteFile(configFile, []byte{}, 0o600))

	path, err := utils.DetermineConfigPath(configFile, true)
	assert.NoError(t, err)
	assert.Equal(t, configFile, path)
}

func Test_Config_Flag_TrailingSeparator_Directory_Rejected(t *testing.T) {
	nonExistentDir := filepath.Join(t.TempDir(), "nonexistent") + string(filepath.Separator)

	_, err := utils.DetermineConfigPath(nonExistentDir, true)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "points to a directory")
	}
}

func Test_Config_EnvVar_TrailingSeparator_Directory_Rejected(t *testing.T) {
	nonExistentDir := filepath.Join(t.TempDir(), "nonexistent") + "/"
	t.Setenv("HARBOR_CLI_CONFIG", nonExistentDir)

	_, err := utils.DetermineConfigPath("", false)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "points to a directory")
	}
}

func Test_InitConfig_Directory_DoesNotCreateOrModifyDataFile(t *testing.T) {
	utils.ConfigInitialization.Reset()
	t.Cleanup(func() {
		utils.ConfigInitialization.Reset()
	})
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "harbor-cli")
	assert.NoError(t, os.MkdirAll(configDir, os.ModePerm))

	dataDir := filepath.Join(tempDir, ".data")
	t.Setenv("XDG_DATA_HOME", dataDir)

	origExitFunc := log.StandardLogger().ExitFunc
	defer func() {
		log.StandardLogger().ExitFunc = origExitFunc
	}()

	var fatalCalled bool
	log.StandardLogger().ExitFunc = func(code int) {
		fatalCalled = true
		panic("log.Fatalf called")
	}

	// 1. Assert that when data.yaml does not exist, calling InitConfig with a directory does not create it.
	assert.Panics(t, func() {
		utils.InitConfig(configDir, true)
	})
	assert.True(t, fatalCalled)

	dataPath := filepath.Join(dataDir, "harbor-cli", "data.yaml")
	assert.NoFileExists(t, dataPath, "data.yaml should not be created on invalid directory config path")

	// 2. Assert that when data.yaml already exists with a valid config, calling InitConfig with a directory does not modify it.
	utils.ConfigInitialization.Reset()
	fatalCalled = false
	assert.NoError(t, os.MkdirAll(filepath.Dir(dataPath), os.ModePerm))
	originalContent := "configpath: /original/valid/config.yaml\n"
	assert.NoError(t, os.WriteFile(dataPath, []byte(originalContent), 0o600))

	assert.Panics(t, func() {
		utils.InitConfig(configDir, true)
	})
	assert.True(t, fatalCalled)

	content, err := os.ReadFile(dataPath)
	assert.NoError(t, err)
	assert.Equal(t, originalContent, string(content), "data.yaml should not be changed on invalid directory config path")
}
