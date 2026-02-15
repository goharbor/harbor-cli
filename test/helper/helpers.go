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
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// HarborTestConfig holds connection details for a Harbor instance used in tests.
type HarborTestConfig struct {
	URL      string
	Username string
	Password string
}

const (
	// Environment variables for local Harbor instance (set by CI workflow)
	EnvHarborURL      = "HARBOR_URL"
	EnvHarborUsername = "HARBOR_USERNAME"
	EnvHarborPassword = "HARBOR_PASSWORD"
	EnvHarborRequired = "HARBOR_REQUIRED"
)

// GetHarborConfig returns the Harbor instance configuration for tests.
// It requires a local Harbor instance (from CI podman setup) via environment variables.
// If Harbor is not available:
// Skips if HARBOR_REQUIRED is not set
// Fails if HARBOR_REQUIRED=true
func GetHarborConfig(t *testing.T) *HarborTestConfig {
	t.Helper()

	localURL := os.Getenv(EnvHarborURL)
	localUsername := os.Getenv(EnvHarborUsername)
	localPassword := os.Getenv(EnvHarborPassword)
	harborRequired := os.Getenv(EnvHarborRequired) == "true"

	if localURL == "" || localUsername == "" || localPassword == "" {
		if harborRequired {
			t.Fatalf("Missing Harbor test environment variables: %s, %s, %s (required for harbor-e2e)", EnvHarborURL, EnvHarborUsername, EnvHarborPassword)
		}
		t.Skipf("Harbor test environment variables not set. Harbor instance required for this test.")
	}

	// If Harbor instance is not healthy
	if !isHarborHealthy(localURL) {
		if harborRequired {
			t.Fatalf("Local Harbor instance at %s is not healthy or unreachable (required for harbor-e2e)", localURL)
		}
		t.Skipf("Local Harbor instance at %s is not healthy or unreachable", localURL)
	}

	t.Logf("Using local Harbor instance at %s", localURL)
	return &HarborTestConfig{
		URL:      localURL,
		Username: localUsername,
		Password: localPassword,
	}
}

// GetHarborServerAddresses returns a list of valid server address formats for the current Harbor instance.
// Useful for testing different URL formats (with/without port, http/https).
func GetHarborServerAddresses(t *testing.T) []string {
	t.Helper()
	cfg := GetHarborConfig(t)

	// Local instance typically runs on a single address
	return []string{cfg.URL}
}

// isHarborHealthy checks if a Harbor instance is responding to health checks.
func isHarborHealthy(baseURL string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // #nosec G402 - Allow self-signed certs for local test instance only
			},
		},
	}

	healthURL := fmt.Sprintf("%s/api/v2.0/health", baseURL)
	resp, err := client.Get(healthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

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
