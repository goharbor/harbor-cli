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
package prompt

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failClient makes every API call fail with errAPI, as it would when Harbor
// is unreachable or the user is not logged in.
func failClient(t *testing.T, errAPI error) {
	t.Helper()

	utils.ClientOnce = sync.Once{}
	utils.ClientOnce.Do(func() {})
	utils.ClientInstance = nil
	utils.ClientErr = errAPI

	t.Cleanup(func() {
		utils.ClientOnce = sync.Once{}
		utils.ClientInstance = nil
		utils.ClientErr = nil
	})
}

func TestPromptsReturnAPIErrors(t *testing.T) {
	errAPI := errors.New("harbor unreachable")

	tests := map[string]func() error{
		"GetRegistryNameFromUser": func() error {
			_, err := GetRegistryNameFromUser()
			return err
		},
		"GetReferenceFromUser": func() error {
			_, err := GetReferenceFromUser("repo", "project")
			return err
		},
		"GetImmutableTagRule": func() error {
			_, err := GetImmutableTagRule("project")
			return err
		},
		"GetTagFromUser": func() error {
			_, err := GetTagFromUser("repo", "project", "sha256:abc")
			return err
		},
		"GetScannerIdFromUser": func() error {
			_, err := GetScannerIdFromUser()
			return err
		},
		"GetMemberIDFromUser": func() error {
			_, err := GetMemberIDFromUser("project", "member")
			return err
		},
	}

	for name, call := range tests {
		t.Run(name, func(t *testing.T) {
			failClient(t, errAPI)

			err := call()
			require.Error(t, err)
			assert.ErrorIs(t, err, errAPI)
		})
	}
}

// emptyListClient points the client at a Harbor where every list is empty.
func emptyListClient(t *testing.T) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]any{})
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

func TestGetRegistryNameFromUserNoRegistries(t *testing.T) {
	emptyListClient(t)

	id, err := GetRegistryNameFromUser()
	require.EqualError(t, err, "no registries found")
	assert.Zero(t, id)
}

func TestPromptsRejectEmptyLists(t *testing.T) {
	tests := map[string]struct {
		call    func() error
		wantErr string
	}{
		"GetReferenceFromUser": {func() error {
			_, err := GetReferenceFromUser("repo", "project")
			return err
		}, "no artifacts found"},
		"GetImmutableTagRule": {func() error {
			_, err := GetImmutableTagRule("project")
			return err
		}, "no immutable tag rules found"},
		"GetTagFromUser": {func() error {
			_, err := GetTagFromUser("repo", "project", "sha256:abc")
			return err
		}, "no tags found"},
		"GetScannerIdFromUser": {func() error {
			_, err := GetScannerIdFromUser()
			return err
		}, "no scanners found"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			emptyListClient(t)
			require.EqualError(t, tt.call(), tt.wantErr)
		})
	}
}
