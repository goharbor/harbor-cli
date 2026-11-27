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
package registry

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestSecretClearedByURLChange(t *testing.T) {
	registry := func(url, secret string) *models.Registry {
		return &models.Registry{
			URL:        url,
			Credential: &models.RegistryCredential{AccessSecret: secret},
		}
	}
	stored := registry("https://old.example.com", "*****")

	tests := []struct {
		name          string
		existing      *models.Registry
		updated       *models.Registry
		secretChanged bool
		want          bool
	}{
		{"url changed without new secret", stored, registry("https://new.example.com", "*****"), false, true},
		{"url changed with new secret", stored, registry("https://new.example.com", "s3cret"), true, false},
		{"url unchanged", stored, registry("https://old.example.com", "*****"), false, false},
		{"only hostname case changed", stored, registry("https://OLD.Example.com", "*****"), false, false},
		{"only trailing slash added", stored, registry("https://old.example.com/", "*****"), false, false},
		{"path case changed", registry("https://old.example.com/Path", "*****"), registry("https://old.example.com/path", "*****"), false, true},
		{"no stored secret", registry("https://old.example.com", ""), registry("https://new.example.com", ""), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, secretClearedByURLChange(tt.existing, tt.updated, tt.secretChanged))
		})
	}
}
