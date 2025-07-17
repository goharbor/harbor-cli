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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/replication"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListReplicationPolicies(opts ...ListFlags) (*replication.ListReplicationPoliciesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags

	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Replication.ListReplicationPolicies(ctx, &replication.ListReplicationPoliciesParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Q:        &listFlags.Q,
		Name:     &listFlags.Name,
		Sort:     &listFlags.Sort,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetReplicationPolicy(policyID int64) (*replication.GetReplicationPolicyOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Replication.GetReplicationPolicy(ctx, &replication.GetReplicationPolicyParams{ID: policyID})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func DeleteReplicationPolicy(policyID int64) (*replication.DeleteReplicationPolicyOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Replication.DeleteReplicationPolicy(ctx, &replication.DeleteReplicationPolicyParams{ID: policyID})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func CreateReplicationPolicy(policy *replication.CreateReplicationPolicyParams) (*replication.CreateReplicationPolicyCreated, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Replication.CreateReplicationPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func UpdateReplicationPolicy(policyID int64, policy *models.ReplicationPolicy) (*replication.UpdateReplicationPolicyOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Replication.UpdateReplicationPolicy(ctx, &replication.UpdateReplicationPolicyParams{
		Context: ctx,
		ID:      policyID,
		Policy:  policy,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func StartReplication(policyID int64) (*replication.StartReplicationCreated, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Replication.StartReplication(ctx, &replication.StartReplicationParams{
		Context: ctx,
		Execution: &models.StartReplicationExecution{
			PolicyID: policyID,
		},
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func StopReplication(policyID int64) (*replication.StopReplicationOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Replication.StopReplication(ctx, &replication.StopReplicationParams{
		Context: ctx,
		ID:      policyID,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ListReplicationExecutions(policyID int64, opts ...ListFlags) (*replication.ListReplicationExecutionsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags

	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Replication.ListReplicationExecutions(ctx, &replication.ListReplicationExecutionsParams{
		Context:  ctx,
		PolicyID: &policyID,
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Sort:     &listFlags.Sort,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetReplicationExecution(executionID int64) (*replication.GetReplicationExecutionOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Replication.GetReplicationExecution(ctx, &replication.GetReplicationExecutionParams{
		Context: ctx,
		ID:      executionID,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetReplicationLog(executionID int64, taskID int64) (*replication.GetReplicationLogOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Replication.GetReplicationLog(ctx, &replication.GetReplicationLogParams{
		Context: ctx,
		ID:      executionID,
		TaskID:  taskID,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ListReplicationTasks(executionID int64, opts ...ListFlags) (*replication.ListReplicationTasksOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags

	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Replication.ListReplicationTasks(ctx, &replication.ListReplicationTasksParams{
		Context:  ctx,
		ID:       executionID,
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Sort:     &listFlags.Sort,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
