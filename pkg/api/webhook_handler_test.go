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
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/webhook/create"
	"github.com/goharbor/harbor-cli/pkg/views/webhook/edit"
	"github.com/stretchr/testify/assert"
)

// Harbor reads a numeric project_name_or_id as a project ID unless
// X-Is-Resource-Name is true, so a project named "123" would otherwise resolve
// to the project whose ID is 123.
func TestWebhookRequestsTreatProjectAsName(t *testing.T) {
	const (
		policiesPath = "/api/v2.0/projects/123/webhook/policies"
		policyPath   = policiesPath + "/7"
	)

	tests := []struct {
		name     string
		body     string
		call     func() error
		wantPath string
	}{
		{
			name: "list",
			body: "[]",
			call: func() error {
				_, err := ListWebhooks("123")
				return err
			},
			wantPath: policiesPath,
		},
		{
			name: "look up id by name",
			body: `[{"id":7,"name":"audit-events"}]`,
			call: func() error {
				_, err := GetWebhookID("123", "audit-events")
				return err
			},
			wantPath: policiesPath,
		},
		{
			name: "get",
			body: "{}",
			call: func() error {
				_, err := GetWebhook("123", 7)
				return err
			},
			wantPath: policyPath,
		},
		{
			name: "create",
			call: func() error {
				return CreateWebhook(&create.CreateView{ProjectName: "123"})
			},
			wantPath: policiesPath,
		},
		{
			name: "update",
			call: func() error {
				return UpdateWebhook(&edit.EditView{ProjectName: "123", WebhookId: 7})
			},
			wantPath: policyPath,
		},
		{
			name: "delete",
			call: func() error {
				return DeleteWebhook("123", 7)
			},
			wantPath: policyPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mu         sync.Mutex
				gotPath    string
				gotIsName  string
				gotRequest bool
			)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				gotRequest = true
				gotPath = r.URL.EscapedPath()
				gotIsName = r.Header.Get("X-Is-Resource-Name")
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodPost {
					w.WriteHeader(http.StatusCreated)
				}
				_, _ = w.Write([]byte(tt.body))
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

			assert.NoError(t, tt.call())

			mu.Lock()
			defer mu.Unlock()
			assert.True(t, gotRequest)
			assert.Equal(t, tt.wantPath, gotPath)
			assert.Equal(t, "true", gotIsName)
		})
	}
}
