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
package api

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
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// robotServer fakes the Harbor endpoints the robot handlers touch and
// records the robot payloads it receives.
type robotServer struct {
	robots       []*models.Robot
	createdRobot *models.RobotCreate
	updatedRobot *models.Robot
}

func useRobotServer(t *testing.T, state *robotServer) {
	t.Helper()

	var mu sync.Mutex
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v2.0/configurations":
			_, _ = w.Write([]byte(`{"robot_name_prefix": {"value": "robot$", "editable": true}}`))
		case r.URL.Path == "/api/v2.0/projects/library":
			_, _ = w.Write([]byte(`{"project_id": 1, "name": "library"}`))
		case r.URL.Path == "/api/v2.0/projects/1":
			_, _ = w.Write([]byte(`{"project_id": 1, "name": "library"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2.0/robots":
			_ = json.NewEncoder(w).Encode(state.robots)
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2.0/robots":
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			state.createdRobot = &models.RobotCreate{}
			_ = json.Unmarshal(body, state.createdRobot)
			mu.Unlock()
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 7, "name": "robot$ci", "secret": "s3cret"}`))
		case r.Method == http.MethodPut && r.URL.Path == "/api/v2.0/robots/7":
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			state.updatedRobot = &models.Robot{}
			_ = json.Unmarshal(body, state.updatedRobot)
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
}

func somePermissions() []*models.RobotPermission {
	return []*models.RobotPermission{{
		Kind:      "project",
		Namespace: "library",
		Access:    []*models.Access{{Resource: "repository", Action: "pull"}},
	}}
}

func TestCreateRobotPassesPermissionsThrough(t *testing.T) {
	state := &robotServer{}
	useRobotServer(t, state)

	response, err := CreateRobot(robotview.CreateView{
		Name:        "ci",
		Level:       "project",
		Duration:    30,
		Permissions: somePermissions(),
	})
	require.NoError(t, err)
	assert.Equal(t, int64(7), response.Payload.ID)

	require.NotNil(t, state.createdRobot)
	require.Len(t, state.createdRobot.Permissions, 1)
	assert.Equal(t, "library", state.createdRobot.Permissions[0].Namespace)
	require.Len(t, state.createdRobot.Permissions[0].Access, 1)
	assert.Equal(t, "pull", state.createdRobot.Permissions[0].Access[0].Action)
}

func TestUpdateRobotPassesPermissionsThrough(t *testing.T) {
	state := &robotServer{}
	useRobotServer(t, state)

	err := UpdateRobot(&robotview.UpdateView{
		ID:          7,
		Name:        "robot$ci",
		Level:       "project",
		Duration:    30,
		Permissions: somePermissions(),
	})
	require.NoError(t, err)

	require.NotNil(t, state.updatedRobot)
	require.Len(t, state.updatedRobot.Permissions, 1)
	assert.Equal(t, "repository", state.updatedRobot.Permissions[0].Access[0].Resource)
}

func TestGetRobotByName(t *testing.T) {
	tests := []struct {
		name    string
		robots  []*models.Robot
		wantErr string
	}{
		{"found", []*models.Robot{{ID: 7, Name: "robot$ci"}}, ""},
		{"missing", []*models.Robot{}, "no robot found with name"},
		{"ambiguous", []*models.Robot{{ID: 7}, {ID: 8}}, "multiple robots found with name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useRobotServer(t, &robotServer{robots: tt.robots})

			// The prefix is stripped before querying, so both spellings work.
			robot, err := GetRobotByName("robot$ci")
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, int64(7), robot.ID)
		})
	}
}

func TestCheckRobotWithNameExists(t *testing.T) {
	useRobotServer(t, &robotServer{robots: []*models.Robot{{ID: 7, Name: "robot$library+ci"}}})

	exists, err := CheckRobotWithNameExists(1, "ci")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = CheckRobotWithNameExists(1, "other")
	require.NoError(t, err)
	assert.False(t, exists)
}
