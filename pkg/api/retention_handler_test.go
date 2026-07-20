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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveRetentionRuleByID(t *testing.T) {
	policy := &models.RetentionPolicy{
		Rules: []*models.RetentionRule{
			{ID: 30},
			nil,
			{ID: 10},
			{ID: 20},
		},
	}

	err := removeRetentionRule(policy, 10)

	require.NoError(t, err)
	require.Len(t, policy.Rules, 3)
	assert.Equal(t, int64(30), policy.Rules[0].ID)
	assert.Nil(t, policy.Rules[1])
	assert.Equal(t, int64(20), policy.Rules[2].ID)
}

func TestRemoveRetentionRuleMissingID(t *testing.T) {
	policy := &models.RetentionPolicy{
		Rules: []*models.RetentionRule{{ID: 10}, {ID: 20}},
	}

	err := removeRetentionRule(policy, 99)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "retention rule with ID 99 not found")
	require.Len(t, policy.Rules, 2)
	assert.Equal(t, int64(10), policy.Rules[0].ID)
	assert.Equal(t, int64(20), policy.Rules[1].ID)
}
