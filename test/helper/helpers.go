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
package helpers

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

func SafeSetEnv(key, value string) {
	envMutex.Lock()
	defer envMutex.Unlock()
	os.Setenv(key, value)
}

func SafeUnsetEnv(key string) {
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
	SafeUnsetEnv("HARBOR_CLI_CONFIG")
	SafeUnsetEnv("XDG_CONFIG_HOME")
	SafeUnsetEnv("XDG_DATA_HOME")
	data = nil
}

func Initialize(t *testing.T, tempDir string) *utils.HarborData {
	utils.ConfigInitialization.Reset() // Reset sync.Once for the test
	SafeSetEnv("XDG_DATA_HOME", filepath.Join(tempDir, ".data"))
	utils.InitConfig(filepath.Join(tempDir, ".config", "config.yaml"), true)
	cds := root.RootCmd()
	err := cds.Execute()
	assert.NoError(t, err, "Expected no error for Root command")
	assert.NoError(t, err, "Expected no error for Root command execution")

	currentData, err := utils.GetCurrentHarborData()
	assert.NoError(t, err, "Expected no error when fetching HarborData")
	return currentData
}

func SetMockKeyring(t *testing.T) {
	mockKeyring := utils.NewMockKeyring()
	utils.SetKeyringProvider(mockKeyring)

	t.Cleanup(func() {
		utils.SetKeyringProvider(&utils.SystemKeyring{})
	})
}
