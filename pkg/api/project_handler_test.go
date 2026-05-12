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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goharbor/harbor-cli/pkg/views/project/create"
)

// TestCreateProject_Validation tests the input validation logic for project creation.
func TestCreateProject_Validation(t *testing.T) {
	t.Run("Invalid Registry ID", func(t *testing.T) {
		opts := create.CreateView{
			ProxyCache: true,
			RegistryID: "abc",
		}
		err := CreateProject(opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid registry ID \"abc\"")
	})

	t.Run("Invalid Storage Limit", func(t *testing.T) {
		opts := create.CreateView{
			StorageLimit: "invalid-limit",
		}
		err := CreateProject(opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid storage format")
	})

	t.Run("Valid Inputs (Validation Pass)", func(t *testing.T) {
		opts := create.CreateView{
			ProxyCache:   true,
			RegistryID:   "123",
			StorageLimit: "10",
		}
		err := CreateProject(opts)
		// We expect an error because we are not logged in, but it should NOT be a validation error
		assert.Error(t, err)
		assert.NotContains(t, err.Error(), "invalid registry ID")
		assert.NotContains(t, err.Error(), "invalid storage format")
	})
}
