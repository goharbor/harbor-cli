package e2e

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var envMutex sync.Mutex

func safeSetEnv(key, value string) {
	envMutex.Lock()
	defer envMutex.Unlock()
	os.Setenv(key, value)
}

func safeUnsetEnv(key string) {
	envMutex.Lock()
	defer envMutex.Unlock()
	os.Unsetenv(key)
}

func ConfigCleanup(t *testing.T, data *utils.HarborData) {
	if data != nil && data.ConfigPath != "" {
		err := os.Remove(data.ConfigPath)
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("Failed to clean up test config file: %v", err)
		}
	}
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		err := os.RemoveAll(os.Getenv("XDG_CONFIG_HOME"))
		if err != nil {
			t.Fatalf("Failed to clean up test config directory: %v", err)
		}
	}
	if os.Getenv("XDG_DATA_HOME") != "" {
		err := os.RemoveAll(os.Getenv("XDG_DATA_HOME"))
		if err != nil {
			t.Fatalf("Failed to clean up test data directory: %v", err)
		}
	}
	safeUnsetEnv("HARBOR_CLI_CONFIG")
	safeUnsetEnv("XDG_CONFIG_HOME")
	safeUnsetEnv("XDG_DATA_HOME")
	data = nil
}

func Initialize(t *testing.T, tempDir string) *utils.HarborData {
	utils.ConfigInitialization.Reset() // Reset sync.Once for the test
	safeSetEnv("XDG_DATA_HOME", filepath.Join(tempDir, ".data"))
	utils.InitConfig(filepath.Join(tempDir, ".config", "config.yaml"), true)
	cds := root.RootCmd()
	err := cds.Execute()
	assert.NoError(t, err, "Expected no error for Root command")
	assert.NoError(t, err, "Expected no error for Root command execution")

	currentData, err := utils.GetCurrentHarborData()
	assert.NoError(t, err, "Expected no error when fetching HarborData")
	return currentData
}

func Test_Config_EnvVar(t *testing.T) {
	utils.ConfigInitialization.Reset() // Reset sync.Once for the test
	tempDir := t.TempDir()
	safeSetEnv("HARBOR_CLI_CONFIG", filepath.Join(tempDir, "config.yaml"))
	safeSetEnv("XDG_DATA_HOME", filepath.Join(tempDir, ".data"))
	utils.InitConfig("", false)
	cds := root.RootCmd()
	err := cds.Execute()
	assert.NoError(t, err, "Expected no error for Root command")
	assert.NoError(t, err, "Expected no error for Root command execution")

	currentData, err := utils.GetCurrentHarborData()
	assert.NoError(t, err, "Expected no error when fetching HarborData")
	defer ConfigCleanup(t, currentData)

	currentConfig, err := utils.GetCurrentHarborConfig()
	assert.NoError(t, err, "Expected no error when fetching HarborConfig")
	assert.NotNil(t, currentConfig, "Configuration should not be nil")
	assert.NotNil(t, currentConfig.CurrentCredentialName, "CurrentCredentialName should not be nil")
	assert.NotNil(t, currentConfig.Credentials, "Credentials should not be nil")
	assert.NotNil(t, currentData.ConfigPath, "ConfigPath should not be nil")
}

func Test_Config_Vanilla(t *testing.T) {
	utils.ConfigInitialization.Reset() // Reset sync.Once for the test
	utils.InitConfig("", false)
	cds := root.RootCmd()
	err := cds.Execute()
	assert.NoError(t, err, "Expected no error for Root command")
	assert.NoError(t, err, "Expected no error for Root command execution")
	currentData, err := utils.GetCurrentHarborData()
	assert.NoError(t, err, "Expected no error when fetching HarborData")
	defer ConfigCleanup(t, currentData)

	currentConfig, err := utils.GetCurrentHarborConfig()
	assert.NoError(t, err, "Expected no error when fetching HarborConfig")
	assert.NotNil(t, currentConfig, "Configuration should not be nil")
	assert.NotNil(t, currentConfig.CurrentCredentialName, "CurrentCredentialName should not be nil")
	assert.NotNil(t, currentConfig.Credentials, "Credentials should not be nil")
	assert.NotNil(t, currentData.ConfigPath, "ConfigPath should not be nil")
}

func Test_Config_Xdg(t *testing.T) {
	utils.ConfigInitialization.Reset() // Reset sync.Once for the test
	tempDir := t.TempDir()
	safeSetEnv("HARBOR_CLI_CONFIG", filepath.Join(tempDir, "config.yaml"))
	safeSetEnv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))
	safeSetEnv("XDG_DATA_HOME", filepath.Join(tempDir, ".data"))
	utils.InitConfig("", false)
	cds := root.RootCmd()
	err := cds.Execute()
	assert.NoError(t, err, "Expected no error for Root command")
	assert.NoError(t, err, "Expected no error for Root command execution")

	currentData, err := utils.GetCurrentHarborData()
	assert.NoError(t, err, "Expected no error when fetching HarborData")
	defer ConfigCleanup(t, currentData)

	currentConfig, err := utils.GetCurrentHarborConfig()
	assert.NoError(t, err, "Expected no error when fetching HarborConfig")
	assert.NotNil(t, currentConfig, "Configuration should not be nil")
	assert.NotNil(t, currentConfig.CurrentCredentialName, "CurrentCredentialName should not be nil")
	assert.NotNil(t, currentConfig.Credentials, "Credentials should not be nil")
	assert.NotNil(t, currentData.ConfigPath, "ConfigPath should not be nil")
}

func Test_Config_Flag(t *testing.T) {
	tempDir := t.TempDir()
	data := Initialize(t, tempDir)
	defer ConfigCleanup(t, data)

	testConfigFile := filepath.Join(tempDir, "config.yaml")
	utils.InitConfig(testConfigFile, true)
	currentConfig, err := utils.GetCurrentHarborConfig()
	assert.NoError(t, err, "Expected no error when fetching HarborConfig")
	assert.NotNil(t, currentConfig, "Configuration should not be nil")
	assert.NotNil(t, currentConfig.CurrentCredentialName, "CurrentCredentialName should not be nil")
	assert.NotNil(t, currentConfig.Credentials, "Credentials should not be nil")
	assert.NotNil(t, data.ConfigPath, "ConfigPath should not be nil")
}
