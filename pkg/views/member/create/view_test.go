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

package create

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleSelectOptions(t *testing.T) {
	opts := RoleSelectOptions()
	assert.Len(t, opts, len(Roles))

	for i, r := range Roles {
		assert.Equal(t, r.Name, opts[i].Key)
		assert.Equal(t, r.ID, opts[i].Value)
	}
}

func TestResolveRoleName(t *testing.T) {
	tests := []struct {
		input  string
		wantID int
	}{
		{"Project Admin", 1},
		{"project_admin", 1},
		{"Admin", 1},
		{"admin", 1},
		{"Developer", 2},
		{"developer", 2},
		{"Guest", 3},
		{"guest", 3},
		{"Maintainer", 4},
		{"maintainer", 4},
		{"Limited Guest", 5},
		{"limited_guest", 5},
		{"limitedguest", 5},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			id, err := ResolveRoleName(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantID, id)
		})
	}

	_, err := ResolveRoleName("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid role name")
}

func TestIsValidRoleID(t *testing.T) {
	for id := 1; id <= 5; id++ {
		assert.True(t, IsValidRoleID(id))
	}
	assert.False(t, IsValidRoleID(0))
	assert.False(t, IsValidRoleID(-1))
	assert.False(t, IsValidRoleID(6))
}

func TestAllowedRolesString(t *testing.T) {
	s := AllowedRolesString()
	assert.Contains(t, s, "1 (Project Admin)")
	assert.Contains(t, s, "2 (Developer)")
	assert.Contains(t, s, "3 (Guest)")
	assert.Contains(t, s, "4 (Maintainer)")
	assert.Contains(t, s, "5 (Limited Guest)")
}

