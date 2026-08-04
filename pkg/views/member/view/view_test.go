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
package view

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestMemberColumns(t *testing.T) {
	t.Run("wide returns all columns", func(t *testing.T) {
		got := memberColumns(true)
		assert.Len(t, got, 6)
	})

	t.Run("non-wide hides Role ID and Project ID", func(t *testing.T) {
		got := memberColumns(false)
		assert.Len(t, got, 4)
		for _, col := range got {
			assert.NotEqual(t, "Role ID", col.Title)
			assert.NotEqual(t, "Project ID", col.Title)
		}
	})

	t.Run("non-wide call does not affect later wide calls", func(t *testing.T) {
		memberColumns(false)
		got := memberColumns(true)
		assert.Len(t, got, 6)
	})

	t.Run("package-level columns are not mutated", func(t *testing.T) {
		memberColumns(false)
		assert.Len(t, columns, 6)
	})
}

func TestMemberRow(t *testing.T) {
	member := &models.ProjectMemberEntity{
		ID:         1,
		ProjectID:  2,
		EntityName: "alice",
		EntityType: "u",
		RoleName:   "projectAdmin",
		RoleID:     3,
	}

	t.Run("wide row has all fields", func(t *testing.T) {
		got := memberRow(member, true)
		assert.Len(t, got, 6)
		assert.Equal(t, "1", got[0])
		assert.Equal(t, "alice", got[1])
		assert.Equal(t, "User", got[2])
		assert.Equal(t, "3", got[4])
		assert.Equal(t, "2", got[5])
	})

	t.Run("non-wide row omits role and project IDs", func(t *testing.T) {
		got := memberRow(member, false)
		assert.Len(t, got, 4)
		assert.Equal(t, "1", got[0])
		assert.Equal(t, "alice", got[1])
		assert.Equal(t, "User", got[2])
	})

	t.Run("group entity type is expanded", func(t *testing.T) {
		group := &models.ProjectMemberEntity{EntityType: "g"}
		got := memberRow(group, false)
		assert.Equal(t, "Group", got[2])
	})
}
