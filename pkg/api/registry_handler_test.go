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
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateRegistrySendsSecretOnlyWhenChanged(t *testing.T) {
	tests := []struct {
		name          string
		secret        string
		secretChanged bool
	}{
		{"masked secret from Harbor", "*****", false},
		{"new secret", "s3cret", true},
		{"new secret that looks like the mask", "*****", true},
		{"cleared secret", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := useTestServer(t, "{}")

			err := UpdateRegistry(&models.Registry{
				Name: "hub",
				URL:  "https://hub.docker.com",
				Credential: &models.RegistryCredential{
					Type:         "basic",
					AccessKey:    "user",
					AccessSecret: tt.secret,
				},
			}, 3, tt.secretChanged)
			require.NoError(t, err)

			req := got()
			assert.Equal(t, "/api/v2.0/registries/3", req.Path)
			body := map[string]any{}
			require.NoError(t, json.Unmarshal(req.Body, &body))
			assert.Equal(t, "user", body["access_key"])
			if tt.secretChanged {
				assert.Equal(t, tt.secret, body["access_secret"])
			} else {
				assert.NotContains(t, body, "access_secret")
			}
		})
	}
}
