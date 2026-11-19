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
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/retention/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRetentionForProjectWithoutPolicy(t *testing.T) {
	var created *models.RetentionPolicy
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2.0/projects/library":
			_, _ = w.Write([]byte(`{"project_id": 7, "name": "library", "metadata": {}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2.0/retentions":
			created = &models.RetentionPolicy{}
			if err := json.NewDecoder(r.Body).Decode(created); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		default:
			http.Error(w, "unexpected request", http.StatusBadRequest)
		}
	}))
	defer srv.Close()

	utils.ClientOnce = sync.Once{}
	utils.ClientOnce.Do(func() {})
	utils.ClientInstance = utils.GetClientByConfig(&harbor.ClientSetConfig{URL: srv.URL})
	utils.ClientErr = nil
	t.Cleanup(func() {
		utils.ClientOnce = sync.Once{}
		utils.ClientInstance = nil
		utils.ClientErr = nil
	})

	_, err := GetRetentionId("library", true)
	require.ErrorIs(t, err, ErrNoRetentionPolicy)

	opts := create.CreateView{
		Action:    "retain",
		Template:  "always",
		Algorithm: "or",
		Scope:     create.RetentionPolicyScope{Level: "project"},
	}
	require.NoError(t, CreateRetention(opts, "library", true))
	require.NotNil(t, created, "expected a retention policy to be created")
	assert.Equal(t, int64(7), created.Scope.Ref)
	require.Len(t, created.Rules, 1)
}
