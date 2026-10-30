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

package member

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveRole_ValidRoleNames(t *testing.T) {
	tests := []struct {
		roleName string
		wantID   int
	}{
		{"Project Admin", 1},
		{"Project_Admin", 1},
		{"project_admin", 1},
		{"projectadmin", 1},
		{"Admin", 1},
		{"admin", 1},
		{"Developer", 2},
		{"developer", 2},
		{"Guest", 3},
		{"guest", 3},
		{"Maintainer", 4},
		{"maintainer", 4},
		{"Limited Guest", 5},
		{"Limited_Guest", 5},
		{"limited_guest", 5},
		{"limitedguest", 5},
	}

	for _, tt := range tests {
		t.Run(tt.roleName, func(t *testing.T) {
			got, err := resolveRole(tt.roleName, 0)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantID, got)
		})
	}
}

func TestResolveRole_ValidRoleIDs(t *testing.T) {
	for id := 1; id <= 5; id++ {
		got, err := resolveRole("", id)
		assert.NoError(t, err)
		assert.Equal(t, id, got)
	}
}

func TestResolveRole_InvalidRoleName(t *testing.T) {
	invalidRoles := []string{"SuperAdmin", "root", "viewer", "Owner", "unknown"}
	for _, role := range invalidRoles {
		t.Run(role, func(t *testing.T) {
			_, err := resolveRole(role, 0)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid role name")
		})
	}
}

func TestResolveRole_InvalidRoleID(t *testing.T) {
	invalidIDs := []int{-1, 6, 10, 99}
	for _, id := range invalidIDs {
		_, err := resolveRole("", id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role ID")
	}
}

func TestResolveRole_MatchingFlags(t *testing.T) {
	tests := []struct {
		roleName string
		roleID   int
		wantID   int
	}{
		{"Guest", 3, 3},
		{"Developer", 2, 2},
		{"Project Admin", 1, 1},
		{"Maintainer", 4, 4},
		{"Limited Guest", 5, 5},
	}
	for _, tt := range tests {
		got, err := resolveRole(tt.roleName, tt.roleID)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantID, got)
	}
}

func TestResolveRole_ConflictingFlags(t *testing.T) {
	tests := []struct {
		roleName string
		roleID   int
	}{
		{"Guest", 1},
		{"Developer", 4},
		{"Project_Admin", 2},
		{"Limited_Guest", 3},
	}
	for _, tt := range tests {
		_, err := resolveRole(tt.roleName, tt.roleID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflicting role flags")
	}
}

func TestResolveRole_EmptyRole(t *testing.T) {
	got, err := resolveRole("", 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, got)
}

func TestCreateMemberCommand_Flags(t *testing.T) {
	cmd := CreateMemberCommand()
	assert.NotNil(t, cmd)

	roleFlag := cmd.Flag("role")
	assert.NotNil(t, roleFlag)
	assert.Equal(t, "", roleFlag.DefValue)

	roleIDFlag := cmd.Flag("roleid")
	assert.NotNil(t, roleIDFlag)
	assert.Equal(t, "0", roleIDFlag.DefValue)

	usernameFlag := cmd.Flag("username")
	assert.NotNil(t, usernameFlag)
	assert.Equal(t, "", usernameFlag.DefValue)
}
