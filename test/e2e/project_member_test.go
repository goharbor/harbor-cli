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

//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Harbor project role IDs.
const (
	roleDeveloper = 2
	roleGuest     = 3
)

// projectMember is the subset of Harbor's project member entity the tests need.
type projectMember struct {
	ID         int64  `json:"id"`
	EntityID   int64  `json:"entity_id"`
	EntityName string `json:"entity_name"`
	RoleID     int64  `json:"role_id"`
}

// createUser creates a local Harbor user, registers its removal with
// t.Cleanup and returns its user ID.
func createUser(t *testing.T, username string) int64 {
	t.Helper()

	doJSON(t, http.MethodPost, "/api/v2.0/users", map[string]any{
		"username": username,
		"email":    username + "@example.com",
		"realname": username,
		"password": "Harbor12345",
	}, nil)

	var users []struct {
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
	}
	doJSON(t, http.MethodGet, "/api/v2.0/users/search?username="+url.QueryEscape(username), nil, &users)

	var id int64
	for _, u := range users {
		if u.Username == username {
			id = u.UserID
		}
	}
	require.NotZero(t, id, "user %s not found after creation", username)

	t.Cleanup(func() {
		doJSON(t, http.MethodDelete, fmt.Sprintf("/api/v2.0/users/%d", id), nil, nil)
	})
	return id
}

func addMember(t *testing.T, project, username string, roleID int64) {
	t.Helper()

	doJSON(t, http.MethodPost, "/api/v2.0/projects/"+project+"/members", map[string]any{
		"role_id":     roleID,
		"member_user": map[string]string{"username": username},
	}, nil)
}

// members returns the project's members keyed by entity name.
func members(t *testing.T, project string) map[string]projectMember {
	t.Helper()

	var list []projectMember
	doJSON(t, http.MethodGet, "/api/v2.0/projects/"+project+"/members?page_size=100", nil, &list)

	byName := make(map[string]projectMember, len(list))
	for _, m := range list {
		byName[m.EntityName] = m
	}
	return byName
}

// Test_ProjectMemberUpdate_ByMemberName checks that `--member <name>` updates
// that user's membership. The member API addresses a membership by its own ID,
// which is not the user's ID.
func Test_ProjectMemberUpdate_ByMemberName(t *testing.T) {
	suffix := time.Now().UnixNano()
	project := fmt.Sprintf("e2e-1159-%d", suffix)
	alice := fmt.Sprintf("e2e-1159-alice-%d", suffix)
	bob := fmt.Sprintf("e2e-1159-bob-%d", suffix)

	aliceUserID := createUser(t, alice)
	createUser(t, bob)
	createProject(t, project)

	addMember(t, project, bob, roleGuest)
	addMember(t, project, alice, roleGuest)

	// The bug only shows when the membership ID differs from the user ID.
	// Re-adding alice gives her membership a fresh ID.
	if members(t, project)[alice].ID == aliceUserID {
		doJSON(t, http.MethodDelete, fmt.Sprintf("/api/v2.0/projects/%s/members/%d", project, members(t, project)[alice].ID), nil, nil)
		addMember(t, project, alice, roleGuest)
	}
	before := members(t, project)
	require.NotEqual(t, aliceUserID, before[alice].ID, "membership ID must differ from user ID")

	loginAndReload(t)

	assert.NoError(t, harbor(t, "project", "member", "update", project,
		"--member", alice, "--roleid", fmt.Sprint(roleDeveloper)))

	after := members(t, project)
	assert.Equal(t, int64(roleDeveloper), after[alice].RoleID, "alice's role should be updated")
	assert.Equal(t, before[bob].RoleID, after[bob].RoleID, "bob's role should not change")
	assert.Equal(t, before[Username()].RoleID, after[Username()].RoleID, "admin's role should not change")
}
