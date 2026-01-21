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
package api_test

import (
	"testing"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/stretchr/testify/assert"
)

// TestListFlags validates the ListFlags structure for proper pagination defaults
func TestListFlags(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		var flags api.ListFlags
		
		assert.Equal(t, int64(0), flags.ProjectID)
		assert.Equal(t, int64(0), flags.Page)
		assert.Equal(t, int64(0), flags.PageSize)
		assert.Equal(t, "", flags.Scope)
		assert.Equal(t, "", flags.Name)
		assert.Equal(t, "", flags.Q)
		assert.Equal(t, "", flags.Sort)
		assert.False(t, flags.Public)
	})

	t.Run("custom pagination", func(t *testing.T) {
		flags := api.ListFlags{
			ProjectID: 10,
			Page:      1,
			PageSize:  50,
			Scope:     "project",
			Name:      "test",
			Q:         "admin",
			Sort:      "name",
			Public:    true,
		}
		
		assert.Equal(t, int64(10), flags.ProjectID)
		assert.Equal(t, int64(1), flags.Page)
		assert.Equal(t, int64(50), flags.PageSize)
		assert.Equal(t, "project", flags.Scope)
		assert.Equal(t, "test", flags.Name)
		assert.Equal(t, "admin", flags.Q)
		assert.Equal(t, "name", flags.Sort)
		assert.True(t, flags.Public)
	})
}

// TestListMemberOptions validates member listing options structure
func TestListMemberOptions(t *testing.T) {
	t.Run("valid member options", func(t *testing.T) {
		opts := api.ListMemberOptions{
			XIsResourceName: true,
			ProjectNameOrID: "library",
			EntityName:      "testuser",
			Page:            1,
			PageSize:        10,
		}
		
		assert.True(t, opts.XIsResourceName)
		assert.Equal(t, "library", opts.ProjectNameOrID)
		assert.Equal(t, "testuser", opts.EntityName)
		assert.Equal(t, int64(1), opts.Page)
		assert.Equal(t, int64(10), opts.PageSize)
	})

	t.Run("project by ID vs name", func(t *testing.T) {
		byName := api.ListMemberOptions{
			XIsResourceName: true,
			ProjectNameOrID: "my-project",
		}
		
		byID := api.ListMemberOptions{
			XIsResourceName: false,
			ProjectNameOrID: "123",
		}
		
		assert.True(t, byName.XIsResourceName)
		assert.False(t, byID.XIsResourceName)
	})
}

// TestUpdateMemberOptions validates member update options structure
func TestUpdateMemberOptions(t *testing.T) {
	t.Run("valid update options", func(t *testing.T) {
		opts := api.UpdateMemberOptions{
			XIsResourceName: true,
			ProjectNameOrID: "library",
			ID:              123,
			RoleID:          nil,
		}
		
		assert.True(t, opts.XIsResourceName)
		assert.Equal(t, "library", opts.ProjectNameOrID)
		assert.Equal(t, int64(123), opts.ID)
		assert.Nil(t, opts.RoleID)
	})
}

// TestGetMemberOptions validates member retrieval options structure
func TestGetMemberOptions(t *testing.T) {
	t.Run("get member by ID", func(t *testing.T) {
		opts := api.GetMemberOptions{
			XIsResourceName: false,
			ProjectNameOrID: "42",
			ID:              100,
		}
		
		assert.False(t, opts.XIsResourceName)
		assert.Equal(t, "42", opts.ProjectNameOrID)
		assert.Equal(t, int64(100), opts.ID)
	})

	t.Run("get member by name", func(t *testing.T) {
		opts := api.GetMemberOptions{
			XIsResourceName: true,
			ProjectNameOrID: "test-project",
			ID:              200,
		}
		
		assert.True(t, opts.XIsResourceName)
		assert.Equal(t, "test-project", opts.ProjectNameOrID)
		assert.Equal(t, int64(200), opts.ID)
	})
}
