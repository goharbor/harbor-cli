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
package policies

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/replication/policies/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToPolicy_DestNamespaceMappings(t *testing.T) {
	// Test that DestNamespace and DestNamespaceReplaceCount are correctly converted
	view := &create.CreateView{
		Name:                      "test-policy",
		Description:               "test description",
		Enabled:                   true,
		Override:                  true,
		ReplicateDeletion:         false,
		DestNamespace:             "dest-ns",
		DestNamespaceReplaceCount: 2,
		CopyByChunk:               true,
		TriggerType:               "manual",
	}

	registry := &models.Registry{
		ID:   1,
		Name: "test-registry",
		Type: "harbor",
	}

	policy := ConvertToPolicy(view, registry)

	// Assert DestNamespace is set correctly
	assert.Equal(t, "dest-ns", policy.DestNamespace)

	// Assert DestNamespaceReplaceCount pointer is set correctly
	assert.NotNil(t, policy.DestNamespaceReplaceCount)
	assert.Equal(t, int8(2), *policy.DestNamespaceReplaceCount)

	// Assert CopyByChunk pointer is set correctly
	assert.NotNil(t, policy.CopyByChunk)
	assert.Equal(t, true, *policy.CopyByChunk)
}

func TestConvertToPolicy_DestNamespaceReplaceCountZero(t *testing.T) {
	// Test that DestNamespaceReplaceCount of 0 is handled correctly (pointer should still be set)
	view := &create.CreateView{
		Name:                      "test-policy",
		Description:               "test description",
		Enabled:                   true,
		Override:                  true,
		ReplicateDeletion:         false,
		DestNamespace:             "dest-ns",
		DestNamespaceReplaceCount: 0,
		CopyByChunk:               false,
		TriggerType:               "manual",
	}

	registry := &models.Registry{
		ID:   1,
		Name: "test-registry",
		Type: "harbor",
	}

	policy := ConvertToPolicy(view, registry)

	// Assert DestNamespace is set correctly
	assert.Equal(t, "dest-ns", policy.DestNamespace)

	// Assert DestNamespaceReplaceCount pointer is set (even for 0)
	assert.NotNil(t, policy.DestNamespaceReplaceCount)
	assert.Equal(t, int8(0), *policy.DestNamespaceReplaceCount)
}

func TestConvertToPolicy_EmptyDestNamespace(t *testing.T) {
	// Test that empty DestNamespace is handled correctly
	view := &create.CreateView{
		Name:                      "test-policy",
		Description:               "test description",
		Enabled:                   true,
		Override:                  true,
		ReplicateDeletion:         false,
		DestNamespace:             "",
		DestNamespaceReplaceCount: 0,
		CopyByChunk:               false,
		TriggerType:               "manual",
	}

	registry := &models.Registry{
		ID:   1,
		Name: "test-registry",
		Type: "harbor",
	}

	policy := ConvertToPolicy(view, registry)

	// Assert DestNamespace is empty
	assert.Equal(t, "", policy.DestNamespace)

	// Assert DestNamespaceReplaceCount pointer is set (even for 0)
	assert.NotNil(t, policy.DestNamespaceReplaceCount)
	assert.Equal(t, int8(0), *policy.DestNamespaceReplaceCount)
}

func TestValidateCreateView(t *testing.T) {
	tests := []struct {
		name          string
		view          *create.CreateView
		expectError   bool
		errorContains string
	}{
		{
			name: "valid with default replace count",
			view: &create.CreateView{
				Name:                      "test-policy",
				DestNamespaceReplaceCount: -1,
			},
			expectError: false,
		},
		{
			name: "valid with explicit replace count and namespace",
			view: &create.CreateView{
				Name:                      "test-policy",
				DestNamespace:             "my-ns",
				DestNamespaceReplaceCount: 2,
			},
			expectError: false,
		},
		{
			name: "invalid replace count -2",
			view: &create.CreateView{
				Name:                      "test-policy",
				DestNamespaceReplaceCount: -2,
			},
			expectError:   true,
			errorContains: "must be between -1 and 3",
		},
		{
			name: "invalid replace count 4",
			view: &create.CreateView{
				Name:                      "test-policy",
				DestNamespaceReplaceCount: 4,
			},
			expectError:   true,
			errorContains: "must be between -1 and 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateView(tt.view)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
