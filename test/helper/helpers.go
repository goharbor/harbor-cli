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
	"github.com/spf13/viper"
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

// reset all of our singletons and viper between runs
func resetAll() {
	// reset config loader (declared as sync.Once)
	utils.ConfigInitialization = &utils.Once{}
	utils.CurrentHarborConfig = nil
	utils.CurrentHarborData = nil

	// reset client singleton (also declared as sync.Once)
	utils.ClientOnce = sync.Once{}
	utils.ClientInstance = nil
	utils.ClientErr = nil

	// wipe viper’s global state
	viper.Reset()
}

func Initialize(t *testing.T, tempDir string) *utils.HarborData {
	resetAll()

	// point the CLI at our test config dir
	cfgDir := filepath.Join(tempDir, ".config")
	dataDir := filepath.Join(tempDir, ".data")
	os.Setenv("XDG_CONFIG_HOME", cfgDir)
	os.Setenv("XDG_DATA_HOME", dataDir)

	// this will create both the data.yaml and the empty config.yaml underneath
	// and return you a *utils.HarborData with the path to config.yaml
	utils.InitConfig(filepath.Join(cfgDir, "config.yaml"), true)
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

// clean up both XDG dirs and reset again
func ConfigCleanup(t *testing.T, data *utils.HarborData) {
	os.RemoveAll(os.Getenv("XDG_CONFIG_HOME"))
	os.RemoveAll(os.Getenv("XDG_DATA_HOME"))
	os.Unsetenv("XDG_CONFIG_HOME")
	os.Unsetenv("XDG_DATA_HOME")
	resetAll()
}

func GetTestHarborURL() string {
	if url := os.Getenv("TEST_HARBOR_URL"); url != "" {
		return url
	}
	return "demo.goharbor.io"
}

func GetTestHarborCredentials() (username, password string) {
	username = os.Getenv("TEST_HARBOR_USERNAME")
	password = os.Getenv("TEST_HARBOR_PASSWORD")

	if username == "" {
		username = "harbor-cli"
	}
	if password == "" {
		password = "Harbor12345"
	}

	return username, password
}
