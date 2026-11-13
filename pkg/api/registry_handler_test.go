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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateRegistryDoesNotSendMaskedSecret(t *testing.T) {
	tests := []struct {
		name       string
		secret     string
		wantSecret any
	}{
		{"masked secret from Harbor", "*****", nil},
		{"new secret", "s3cret", "s3cret"},
		{"cleared secret", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mu sync.Mutex
			body := map[string]any{}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				raw, _ := io.ReadAll(r.Body)
				mu.Lock()
				_ = json.Unmarshal(raw, &body)
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte("{}"))
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

			err := UpdateRegistry(&models.Registry{
				Name: "hub",
				URL:  "https://hub.docker.com",
				Credential: &models.RegistryCredential{
					Type:         "basic",
					AccessKey:    "user",
					AccessSecret: tt.secret,
				},
			}, 3)
			require.NoError(t, err)

			mu.Lock()
			defer mu.Unlock()
			assert.Equal(t, "user", body["access_key"])
			assert.Equal(t, tt.wantSecret, body["access_secret"])
		})
	}
}
