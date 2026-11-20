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
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testRobot = `{
	"id": 7,
	"name": "robot$ci",
	"level": "system",
	"description": "old description",
	"duration": 30,
	"permissions": [
		{"kind": "system", "namespace": "/", "access": [{"resource": "registry", "action": "list"}]}
	]
}`

type robotUpdate struct {
	path  string
	robot models.Robot
}

const testRobotWithoutPermissions = `{
	"id": 7,
	"name": "robot$ci",
	"level": "system",
	"description": "old description",
	"duration": 30,
	"permissions": []
}`

func useRobotServer(t *testing.T) *robotUpdate {
	t.Helper()
	return useRobotServerWith(t, testRobot)
}

func useRobotServerWith(t *testing.T, robotJSON string) *robotUpdate {
	t.Helper()

	var mu sync.Mutex
	updated := &robotUpdate{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2.0/robots/7":
			_, _ = w.Write([]byte(robotJSON))
		case r.Method == http.MethodPut:
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			updated.path = r.URL.Path
			_ = json.Unmarshal(body, &updated.robot)
			mu.Unlock()
			_, _ = w.Write([]byte("{}"))
		default:
			http.Error(w, "unexpected request", http.StatusNotFound)
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

	return updated
}

func TestUpdateRobotCommandUsesFlags(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantDescription string
		wantDuration    int64
	}{
		{"description", []string{"7", "--description", "nightly builds"}, "nightly builds", 30},
		{"duration", []string{"7", "--duration", "90"}, "old description", 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated := useRobotServer(t)

			cmd := UpdateRobotCommand()
			cmd.SetArgs(tt.args)
			require.NoError(t, cmd.Execute())

			assert.Equal(t, "/api/v2.0/robots/7", updated.path)
			assert.Equal(t, tt.wantDescription, updated.robot.Description)
			require.NotNil(t, updated.robot.Duration)
			assert.Equal(t, tt.wantDuration, *updated.robot.Duration)
			require.Len(t, updated.robot.Permissions, 1)
			assert.Equal(t, "system", updated.robot.Permissions[0].Kind)
			assert.Equal(t, "/", updated.robot.Permissions[0].Namespace)
			require.Len(t, updated.robot.Permissions[0].Access, 1)
			assert.Equal(t, "registry", updated.robot.Permissions[0].Access[0].Resource)
			assert.Equal(t, "list", updated.robot.Permissions[0].Access[0].Action)
		})
	}
}

func TestUpdateRobotCommandRejectsZeroDuration(t *testing.T) {
	updated := useRobotServer(t)

	cmd := UpdateRobotCommand()
	cmd.SetArgs([]string{"7", "--duration", "0"})
	require.ErrorContains(t, cmd.Execute(), "duration cannot be 0")
	assert.Empty(t, updated.path)
}

// askForPermissions replaces the permission prompts for the length of a test.
func askForPermissions(t *testing.T, chosen []models.Permission) *bool {
	t.Helper()
	asked := false
	original := askForPermissionsForUpdate
	askForPermissionsForUpdate = func(_ *update.UpdateView, permissions *[]models.Permission, _ map[string][]models.Permission) error {
		asked = true
		*permissions = append(*permissions, chosen...)
		return nil
	}
	t.Cleanup(func() { askForPermissionsForUpdate = original })
	return &asked
}

func TestUpdateRobotCommandAsksForPermissionsWhenTheRobotHasNone(t *testing.T) {
	updated := useRobotServerWith(t, testRobotWithoutPermissions)
	asked := askForPermissions(t, []models.Permission{{Resource: "repository", Action: "pull"}})

	cmd := UpdateRobotCommand()
	cmd.SetArgs([]string{"7", "--description", "nightly builds"})
	require.NoError(t, cmd.Execute())

	assert.True(t, *asked)
	assert.Equal(t, "nightly builds", updated.robot.Description)
	require.Len(t, updated.robot.Permissions, 1)
	assert.Equal(t, "system", updated.robot.Permissions[0].Kind)
	require.Len(t, updated.robot.Permissions[0].Access, 1)
	assert.Equal(t, "repository", updated.robot.Permissions[0].Access[0].Resource)
	assert.Equal(t, "pull", updated.robot.Permissions[0].Access[0].Action)
}

func TestUpdateRobotCommandKeepsPermissionsWithoutAsking(t *testing.T) {
	updated := useRobotServer(t)
	asked := askForPermissions(t, nil)

	cmd := UpdateRobotCommand()
	cmd.SetArgs([]string{"7", "--description", "nightly builds"})
	require.NoError(t, cmd.Execute())

	assert.False(t, *asked)
	require.Len(t, updated.robot.Permissions, 1)
}

func TestUpdateRobotCommandRejectsARobotLeftWithoutPermissions(t *testing.T) {
	updated := useRobotServerWith(t, testRobotWithoutPermissions)
	asked := askForPermissions(t, nil)

	cmd := UpdateRobotCommand()
	cmd.SetArgs([]string{"7", "--description", "nightly builds"})
	require.ErrorContains(t, cmd.Execute(), "robot account needs at least one permission")

	assert.True(t, *asked)
	assert.Empty(t, updated.path)
}
