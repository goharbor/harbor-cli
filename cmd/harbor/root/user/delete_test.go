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
	"net/url"
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

// fakeUserServer is a minimal Harbor user API. It follows Harbor's q syntax
// for usernames: q is URL-decoded once more after the HTTP layer,
// "username=v" is an exact match with one leading backslash removed,
// "username=~v" is a substring match, and a value starting with [, { or ( is
// rejected as a range or list.
type fakeUserServer struct {
	mu          sync.Mutex
	users       map[string]int64
	forbidden   bool
	unfiltered  int
	listQueries []string
	deletedIDs  []int64
}

func (f *fakeUserServer) hasID(id int64) bool {
	for _, uid := range f.users {
		if uid == id {
			return true
		}
	}
	return false
}

func (f *fakeUserServer) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	f.listQueries = append(f.listQueries, q)
	if q == "" {
		f.unfiltered++
	}

	if unescaped, err := url.QueryUnescape(q); err == nil {
		q = unescaped
	}

	var names []string
	switch v, ok := strings.CutPrefix(q, "username="); {
	case q == "":
		for name := range f.users {
			names = append(names, name)
		}
	case !ok:
		w.WriteHeader(http.StatusBadRequest)
		return
	case strings.HasPrefix(v, "~"):
		for name := range f.users {
			if strings.Contains(name, v[1:]) {
				names = append(names, name)
			}
		}
	case strings.ContainsAny(v[:1], "[{("):
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		v = strings.TrimPrefix(v, `\`)
		if _, found := f.users[v]; found {
			names = append(names, v)
		}
	}
	sort.Strings(names)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 || size < 1 {
		// Harbor answers page=0/page_size=0 with an empty list.
		names = nil
	} else {
		start := min((page-1)*size, len(names))
		names = names[start:min(start+size, len(names))]
	}

	payload := []*models.UserResp{}
	for _, name := range names {
		payload = append(payload, &models.UserResp{Username: name, UserID: f.users[name]})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

// useFakeUserServer points the shared client at a fake Harbor that knows the
// given users and records list queries and deletes.
func useFakeUserServer(t *testing.T, users map[string]int64) *fakeUserServer {
	t.Helper()

	f := &fakeUserServer{users: users}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()

		if f.forbidden {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		idPart, isItem := strings.CutPrefix(r.URL.Path, "/api/v2.0/users/")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2.0/users":
			f.list(w, r)
		case isItem && (r.Method == http.MethodGet || r.Method == http.MethodDelete):
			id, err := strconv.ParseInt(idPart, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !f.hasID(id) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if r.Method == http.MethodDelete {
				f.deletedIDs = append(f.deletedIDs, id)
				w.WriteHeader(http.StatusOK)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(&models.UserResp{UserID: id})
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

func (f *fakeUserServer) deleted() []int64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	ids := append([]int64(nil), f.deletedIDs...)
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func runUserDelete(args ...string) error {
	cmd := UserDeleteCmd()
	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	return cmd.Execute()
}

func TestUserDeleteResolvesNamesWithoutListingAllUsers(t *testing.T) {
	f := useFakeUserServer(t, map[string]int64{"alice": 3, "bob": 4, "carol": 5})

	require.NoError(t, runUserDelete("alice", "bob", "carol"))

	assert.Equal(t, []int64{3, 4, 5}, f.deleted())
	assert.Zero(t, f.unfiltered, "user delete should look names up, not list every user")
}

func TestUserDeleteUnknownNameDeletesNothing(t *testing.T) {
	f := useFakeUserServer(t, map[string]int64{"alice": 3, "bob": 4})

	err := runUserDelete("alice", "missing403", "bob")
	require.Error(t, err)
	assert.Contains(t, err.Error(), `user "missing403" not found`)
	assert.NotContains(t, err.Error(), "permission denied")
	assert.Empty(t, f.deleted(), "no user should be deleted when an argument does not resolve")
}

func TestUserDeleteByNumericID(t *testing.T) {
	f := useFakeUserServer(t, map[string]int64{"alice": 7})

	require.NoError(t, runUserDelete("7"))
	assert.Equal(t, []int64{7}, f.deleted())
}

func TestUserDeleteNumericNameMatchingAnotherIDIsAmbiguous(t *testing.T) {
	// User "5" has ID 9, and a different user has ID 5.
	f := useFakeUserServer(t, map[string]int64{"5": 9, "alice": 5})

	err := runUserDelete("5")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ambiguous")
	assert.Empty(t, f.deleted())
}

func TestUserDeleteNumericNameWithoutMatchingID(t *testing.T) {
	f := useFakeUserServer(t, map[string]int64{"5": 9})

	require.NoError(t, runUserDelete("5"))
	assert.Equal(t, []int64{9}, f.deleted())
}

func TestUserDeleteNameWithQuerySyntax(t *testing.T) {
	f := useFakeUserServer(t, map[string]int64{"[alice]": 3, "alice": 4, "x[alice]y": 5})

	require.NoError(t, runUserDelete("[alice]"))
	assert.Equal(t, []int64{3}, f.deleted())
}

func TestUserDeleteNamesHarborWouldDecode(t *testing.T) {
	users := map[string]int64{
		"a+b":       3,
		"a b":       4,
		`\alice`:    5,
		"alice":     6,
		"{x}":       7,
		"(y)":       8,
		"tom&jerry": 10,
	}
	f := useFakeUserServer(t, users)

	require.NoError(t, runUserDelete("a+b", `\alice`, "{x}", "(y)", "tom&jerry"))
	assert.Equal(t, []int64{3, 5, 7, 8, 10}, f.deleted())
}

func TestUserDeleteLookupForbidden(t *testing.T) {
	f := useFakeUserServer(t, map[string]int64{"alice": 3})
	f.forbidden = true

	err := runUserDelete("alice")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
}

func TestUserPasswordNumericArgIsNameOnly(t *testing.T) {
	useFakeUserServer(t, map[string]int64{"alice": 7})

	cmd := UserPasswordChangeCmd()
	cmd.SetArgs([]string{"7"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), `user "7" not found`)
}
