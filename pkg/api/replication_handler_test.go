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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/replication"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllReplicationPoliciesFetchesAllPages(t *testing.T) {
	var pages []int64
	var pageSizes []int64

	policies, err := GetAllReplicationPolicies(
		func(opts ...ListFlags) (*replication.ListReplicationPoliciesOK, error) {
			require.Len(t, opts, 1)
			pages = append(pages, opts[0].Page)
			pageSizes = append(pageSizes, opts[0].PageSize)

			if opts[0].Page == 1 {
				return &replication.ListReplicationPoliciesOK{Payload: makeReplicationPolicies(100)}, nil
			}
			return &replication.ListReplicationPoliciesOK{Payload: makeReplicationPolicies(3)}, nil
		},
		ListFlags{Page: 5, PageSize: 0},
	)

	require.NoError(t, err)
	assert.Len(t, policies, 103)
	assert.Equal(t, []int64{1, 2}, pages)
	assert.Equal(t, []int64{100, 100}, pageSizes)
}

func TestGetAllReplicationExecutionsFetchesAllPages(t *testing.T) {
	var pages []int64
	var pageSizes []int64
	const policyID int64 = 7

	executions, err := GetAllReplicationExecutions(
		policyID,
		func(gotPolicyID int64, opts ...ListFlags) (*replication.ListReplicationExecutionsOK, error) {
			require.Equal(t, policyID, gotPolicyID)
			require.Len(t, opts, 1)
			pages = append(pages, opts[0].Page)
			pageSizes = append(pageSizes, opts[0].PageSize)

			if opts[0].Page == 1 {
				return &replication.ListReplicationExecutionsOK{Payload: makeReplicationExecutions(100)}, nil
			}
			return &replication.ListReplicationExecutionsOK{Payload: makeReplicationExecutions(4)}, nil
		},
		ListFlags{Page: 3, PageSize: 0},
	)

	require.NoError(t, err)
	assert.Len(t, executions, 104)
	assert.Equal(t, []int64{1, 2}, pages)
	assert.Equal(t, []int64{100, 100}, pageSizes)
}

func makeReplicationPolicies(count int) []*models.ReplicationPolicy {
	policies := make([]*models.ReplicationPolicy, count)
	for i := range policies {
		policies[i] = &models.ReplicationPolicy{ID: int64(i + 1)}
	}
	return policies
}

func makeReplicationExecutions(count int) []*models.ReplicationExecution {
	executions := make([]*models.ReplicationExecution, count)
	for i := range executions {
		executions[i] = &models.ReplicationExecution{ID: int64(i + 1)}
	}
	return executions
}
