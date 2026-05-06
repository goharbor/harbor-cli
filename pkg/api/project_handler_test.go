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
			StorageLimit: "10GiB",
		}
		err := CreateProject(opts)
		// We expect an error because we are not logged in, but it should NOT be a validation error
		assert.Error(t, err)
		assert.NotContains(t, err.Error(), "invalid registry ID")
		assert.NotContains(t, err.Error(), "invalid storage format")
	})
}
