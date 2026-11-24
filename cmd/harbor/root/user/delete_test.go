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
package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUserServer struct {
	mu         sync.Mutex
	listCalls  int
	deletedIDs []int64
}

// useFakeUserServer points the shared client at a fake Harbor that knows the
// given users and records list and delete calls.
func useFakeUserServer(t *testing.T, users map[string]int64) *fakeUserServer {
	t.Helper()

	f := &fakeUserServer{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2.0/users":
			f.listCalls++
			payload := []*models.UserResp{}
			for name, id := range users {
				payload = append(payload, &models.UserResp{Username: name, UserID: id})
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(payload)
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/v2.0/users/"):
			id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/v2.0/users/"), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			f.deletedIDs = append(f.deletedIDs, id)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
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
	return f
}

func TestUserDeleteListsUsersOnce(t *testing.T) {
	f := useFakeUserServer(t, map[string]int64{"alice": 3, "bob": 4, "carol": 5})

	cmd := UserDeleteCmd()
	cmd.SetArgs([]string{"alice", "bob", "carol"})
	require.NoError(t, cmd.Execute())

	f.mu.Lock()
	defer f.mu.Unlock()
	assert.Equal(t, 1, f.listCalls, "user list should be fetched once, not once per username")
	sort.Slice(f.deletedIDs, func(i, j int) bool { return f.deletedIDs[i] < f.deletedIDs[j] })
	assert.Equal(t, []int64{3, 4, 5}, f.deletedIDs)
}
