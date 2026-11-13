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

func useRobotServer(t *testing.T) *models.Robot {
	t.Helper()

	var mu sync.Mutex
	updated := &models.Robot{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPut {
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			_ = json.Unmarshal(body, updated)
			mu.Unlock()
			_, _ = w.Write([]byte("{}"))
			return
		}
		_, _ = w.Write([]byte(testRobot))
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

			assert.Equal(t, tt.wantDescription, updated.Description)
			require.NotNil(t, updated.Duration)
			assert.Equal(t, tt.wantDuration, *updated.Duration)
			require.Len(t, updated.Permissions, 1)
			assert.Equal(t, "system", updated.Permissions[0].Kind)
			require.Len(t, updated.Permissions[0].Access, 1)
			assert.Equal(t, "registry", updated.Permissions[0].Access[0].Resource)
		})
	}
}
