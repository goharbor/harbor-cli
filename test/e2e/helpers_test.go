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

//go:build e2e

package e2e

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

// loginAndReload runs `harbor login` as Username, then reloads the saved
// config the way a new CLI process would, so later commands and direct API
// helper calls pick up the stored credential.
func loginAndReload(t *testing.T) {
	t.Helper()

	require.NoError(t, login(t, Server(), Username(), Password()))
	reloadConfig(t)
}

// reloadConfig drops the cached config and API client and loads the config
// file that `harbor login` wrote. It returns that file's path.
func reloadConfig(t *testing.T) string {
	t.Helper()

	data, err := utils.GetCurrentHarborData()
	require.NoError(t, err)

	utils.ConfigInitialization = &utils.Once{}
	utils.CurrentHarborConfig = nil
	utils.ClientOnce = sync.Once{}
	utils.ClientInstance = nil
	utils.ClientErr = nil
	utils.InitConfig(data.ConfigPath, true)

	return data.ConfigPath
}

// harbor runs the CLI with args as a fresh invocation would, against the
// config that `harbor login` wrote earlier in the test.
func harbor(t *testing.T, args ...string) error {
	t.Helper()

	configPath := reloadConfig(t)

	cmd := root.RootCmd()
	cmd.SetArgs(append([]string{"--config", configPath}, args...))
	return cmd.Execute()
}

// createProject creates a private project and registers its removal with
// t.Cleanup.
func createProject(t *testing.T, name string) {
	t.Helper()

	doJSON(t, http.MethodPost, "/api/v2.0/projects", map[string]any{
		"project_name": name,
		"metadata":     map[string]string{"public": "false"},
	}, nil)

	t.Cleanup(func() {
		doJSON(t, http.MethodDelete, "/api/v2.0/projects/"+url.PathEscape(name), nil, nil)
	})
}

// captureStdout runs fn and returns what it wrote to os.Stdout.
func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = w
	// Buffered so the reader never blocks if we return early.
	done := make(chan []byte, 1)
	go func() {
		out, _ := io.ReadAll(r)
		done <- out
	}()

	// Restore stdout and close the pipe even if fn calls t.FailNow.
	var restoreOnce sync.Once
	restore := func() {
		restoreOnce.Do(func() {
			os.Stdout = orig
			_ = w.Close()
		})
	}
	defer func() {
		restore()
		_ = r.Close()
	}()

	runErr := fn()
	restore()
	out := <-done

	return string(bytes.TrimSpace(out)), runErr
}
