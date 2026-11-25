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
package robot

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const existingRobot = `{
	"id": 7,
	"name": "robot$ci",
	"level": "system",
	"description": "old description",
	"duration": 30,
	"permissions": [
		{"kind": "system", "namespace": "/", "access": [{"resource": "registry", "action": "list"}]}
	]
}`

type recordedRobots struct {
	created *models.RobotCreate
	updated *models.Robot
}

// useHarborServer fakes every endpoint the create and update pipelines
// touch and records the robot payloads they send.
func useHarborServer(t *testing.T) *recordedRobots {
	t.Helper()

	var mu sync.Mutex
	state := &recordedRobots{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v2.0/permissions":
			_, _ = w.Write([]byte(`{
				"system": [{"resource": "registry", "action": "list"}],
				"project": [{"resource": "repository", "action": "pull"}, {"resource": "repository", "action": "push"}]
			}`))
		case r.URL.Path == "/api/v2.0/configurations":
			_, _ = w.Write([]byte(`{"robot_name_prefix": {"value": "robot$", "editable": true}}`))
		case r.URL.Path == "/api/v2.0/projects/library", r.URL.Path == "/api/v2.0/projects/1":
			_, _ = w.Write([]byte(`{"project_id": 1, "name": "library"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2.0/robots/7":
			_, _ = w.Write([]byte(existingRobot))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2.0/robots":
			_, _ = w.Write([]byte("[]"))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2.0/robots":
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			state.created = &models.RobotCreate{}
			_ = json.Unmarshal(body, state.created)
			mu.Unlock()
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 7, "name": "robot$ci", "secret": "s3cret", "expires_at": 123}`))
		case r.Method == http.MethodPut && r.URL.Path == "/api/v2.0/robots/7":
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			state.updated = &models.Robot{}
			_ = json.Unmarshal(body, state.updated)
			mu.Unlock()
			_, _ = w.Write([]byte("{}"))
		default:
			http.Error(w, "unexpected request: "+r.Method+" "+r.URL.Path, http.StatusNotFound)
		}
	}))

	utils.ClientOnce = sync.Once{}
	utils.ClientOnce.Do(func() {})
	utils.ClientInstance = utils.GetClientByConfig(&harbor.ClientSetConfig{URL: srv.URL})
	utils.ClientErr = nil

	t.Cleanup(func() {
		srv.Close()
		utils.ClientOnce = sync.Once{}
		utils.ClientInstance = nil
		utils.ClientErr = nil
	})

	return state
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "robot-config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0600))
	return path
}

func TestRunCreateFromConfigFile(t *testing.T) {
	state := useHarborServer(t)
	t.Chdir(t.TempDir())

	configPath := writeConfig(t, `
name: ci
description: nightly builds
duration: 90
level: project
permissions:
  - namespace: library
    kind: project
    access:
      - resource: repository
        actions: ["pull"]
`)

	err := RunCreate(CreateOptions{
		Level:        LevelProject,
		View:         &robotview.CreateView{},
		ConfigFile:   configPath,
		ExportToFile: true,
	})
	require.NoError(t, err)

	require.NotNil(t, state.created)
	assert.Equal(t, "ci", state.created.Name)
	assert.Equal(t, "project", state.created.Level)
	require.Len(t, state.created.Permissions, 1)
	assert.Equal(t, "library", state.created.Permissions[0].Namespace)

	// --export-to-file wrote the secret instead of opening the TUI
	assert.FileExists(t, "robot$ci-secret.json")
}

func TestRunCreateRejectsWrongLevel(t *testing.T) {
	useHarborServer(t)

	configPath := writeConfig(t, `
name: ci
duration: 90
level: project
permissions:
  - namespace: library
    kind: project
    access:
      - resource: repository
        actions: ["pull"]
`)

	err := RunCreate(CreateOptions{
		Level:      LevelSystem,
		View:       &robotview.CreateView{},
		ConfigFile: configPath,
	})
	require.ErrorContains(t, err, "level must be 'system'")
}

func TestRunUpdateFromConfigFileReplacesPermissions(t *testing.T) {
	state := useHarborServer(t)

	configPath := writeConfig(t, `
name: ci
description: updated description
duration: 60
level: system
permissions:
  - namespace: library
    kind: project
    access:
      - resource: repository
        actions: ["push"]
`)

	err := RunUpdate(UpdateOptions{
		Level:      LevelSystem,
		RobotID:    7,
		ConfigFile: configPath,
	})
	require.NoError(t, err)

	require.NotNil(t, state.updated)
	assert.Equal(t, "updated description", state.updated.Description)
	require.NotNil(t, state.updated.Duration)
	assert.Equal(t, int64(60), *state.updated.Duration)

	// The config file describes the desired state: the robot's previous
	// system scope is replaced by the file's project scope.
	require.Len(t, state.updated.Permissions, 1)
	assert.Equal(t, "project", state.updated.Permissions[0].Kind)
	assert.Equal(t, "push", state.updated.Permissions[0].Access[0].Action)
}
