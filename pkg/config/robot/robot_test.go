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
package config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func usePermissionsServer(t *testing.T) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2.0/permissions":
			_, _ = w.Write([]byte(`{
				"system": [{"resource": "registry", "action": "list"}, {"resource": "registry", "action": "read"}],
				"project": [{"resource": "repository", "action": "pull"}, {"resource": "repository", "action": "push"}]
			}`))
		case "/api/v2.0/projects/library":
			_, _ = w.Write([]byte(`{"project_id": 1, "name": "library"}`))
		default:
			http.Error(w, "unexpected request: "+r.URL.Path, http.StatusNotFound)
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
}

func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0600))
	return path
}

func TestLoadConfigFromFileYAML(t *testing.T) {
	usePermissionsServer(t)

	path := writeFile(t, "robot.yaml", `
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

	view, err := LoadConfigFromFile(path)
	require.NoError(t, err)
	assert.Equal(t, "ci", view.Name)
	assert.Equal(t, "project", view.Level)
	assert.Equal(t, "library", view.ProjectName)
	require.Len(t, view.Permissions, 1)
	require.Len(t, view.Permissions[0].Access, 1)
	assert.Equal(t, "pull", view.Permissions[0].Access[0].Action)
}

func TestLoadConfigFromFileJSON(t *testing.T) {
	usePermissionsServer(t)

	path := writeFile(t, "robot.json", `{
		"name": "ci",
		"duration": 90,
		"level": "system",
		"permissions": [
			{"namespace": "/", "kind": "system", "access": [{"resource": "registry", "actions": ["list"]}]}
		]
	}`)

	view, err := LoadConfigFromFile(path)
	require.NoError(t, err)
	assert.Equal(t, "system", view.Level)
	assert.Equal(t, "list", view.Permissions[0].Access[0].Action)
}

func TestLoadConfigFromFileExpandsWildcards(t *testing.T) {
	usePermissionsServer(t)

	path := writeFile(t, "robot.yaml", `
name: ci
duration: 90
level: project
permissions:
  - namespace: library
    kind: project
    access:
      - resource: repository
        actions: ["*"]
`)

	view, err := LoadConfigFromFile(path)
	require.NoError(t, err)
	assert.Len(t, view.Permissions[0].Access, 2) // pull and push
}

func TestLoadConfigFromFileSkipsInvalidEntries(t *testing.T) {
	usePermissionsServer(t)

	path := writeFile(t, "robot.yaml", `
name: ci
duration: 90
level: project
permissions:
  - namespace: library
    kind: project
    access:
      - resource: repository
        actions: ["pull", "fly"]
      - resource: spaceship
        actions: ["launch"]
`)

	view, err := LoadConfigFromFile(path)
	require.NoError(t, err)
	// The invalid action and resource are dropped with a warning.
	require.Len(t, view.Permissions, 1)
	require.Len(t, view.Permissions[0].Access, 1)
	assert.Equal(t, "pull", view.Permissions[0].Access[0].Action)
}

func TestLoadConfigFromFileRejectsBadSpecs(t *testing.T) {
	usePermissionsServer(t)

	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			"no permissions", `
name: ci
duration: 90
level: system
permissions: []
`, "no permissions defined"},
		{
			"zero duration", `
name: ci
duration: 0
level: system
permissions:
  - namespace: "/"
    kind: system
    access:
      - resource: registry
        actions: ["list"]
`, "duration cannot be 0"},
		{
			"system scope on project robot", `
name: ci
duration: 90
level: project
permissions:
  - namespace: "/"
    kind: system
    access:
      - resource: registry
        actions: ["list"]
`, "project robots can only have project permission scopes"},
		{
			"two scopes on project robot", `
name: ci
duration: 90
level: project
permissions:
  - namespace: library
    kind: project
    access:
      - resource: repository
        actions: ["pull"]
  - namespace: other
    kind: project
    access:
      - resource: repository
        actions: ["push"]
`, "project robots can only have one permission scope"},
		{
			"wildcard namespace", `
name: ci
duration: 90
level: project
permissions:
  - namespace: "*"
    kind: project
    access:
      - resource: repository
        actions: ["pull"]
`, "wildcard namespace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeFile(t, "robot.yaml", tt.content)
			_, err := LoadConfigFromFile(path)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
