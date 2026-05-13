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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/preheat"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListPreheatPolicies(projectName string, isID bool, opts ...ListFlags) (*preheat.ListPoliciesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}

	if isID {
		project, err := GetProject(projectName, true)
		if err != nil {
			return nil, err
		}
		projectName = project.Payload.Name
	}

	response, err := client.Preheat.ListPolicies(ctx, &preheat.ListPoliciesParams{
		ProjectName: projectName,
		Page:        &listFlags.Page,
		PageSize:    &listFlags.PageSize,
		Q:           &listFlags.Q,
		Sort:        &listFlags.Sort,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetPreheatPolicy(projectName, policyName string) (*preheat.GetPolicyOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Preheat.GetPolicy(ctx, &preheat.GetPolicyParams{
		ProjectName:       projectName,
		PreheatPolicyName: policyName,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func DeletePreheatPolicy(projectName, policyName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Preheat.DeletePolicy(ctx, &preheat.DeletePolicyParams{
		ProjectName:       projectName,
		PreheatPolicyName: policyName,
	})
	if err != nil {
		return err
	}
	return nil
}

func StartPreheatPolicy(projectName, policyName string) (string, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return "", err
	}

	policy, err := GetPreheatPolicy(projectName, policyName)
	if err != nil {
		return "", err
	}

	resp, err := client.Preheat.ManualPreheat(ctx, &preheat.ManualPreheatParams{
		ProjectName:       projectName,
		PreheatPolicyName: policyName,
		Policy:            policy.Payload,
	})
	if err != nil {
		return "", err
	}
	return resp.Location, nil
}

func CreatePreheatPolicy(projectName string, policy *models.PreheatPolicy) (*preheat.CreatePolicyCreated, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Preheat.CreatePolicy(ctx, &preheat.CreatePolicyParams{
		ProjectName: projectName,
		Policy:      policy,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func UpdatePreheatPolicy(projectName, policyName string, policy *models.PreheatPolicy) (*preheat.UpdatePolicyOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Preheat.UpdatePolicy(ctx, &preheat.UpdatePolicyParams{
		ProjectName:       projectName,
		PreheatPolicyName: policyName,
		Policy:            policy,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func ListProvidersUnderProject(projectName string) ([]*models.ProviderUnderProject, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Preheat.ListProvidersUnderProject(ctx, &preheat.ListProvidersUnderProjectParams{
		ProjectName: projectName,
	})
	if err != nil {
		return nil, err
	}
	return response.Payload, nil
}

func ListPreheatExecutions(projectName string, policyName string, opts ...ListFlags) (*preheat.ListExecutionsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Preheat.ListExecutions(ctx, &preheat.ListExecutionsParams{
		ProjectName:       projectName,
		PreheatPolicyName: policyName,
		Page:              &listFlags.Page,
		PageSize:          &listFlags.PageSize,
		Q:                 &listFlags.Q,
		Sort:              &listFlags.Sort,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetPreheatExecution(projectName string, policyName string, executionID int64) (*preheat.GetExecutionOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Preheat.GetExecution(ctx, &preheat.GetExecutionParams{
		ProjectName:       projectName,
		PreheatPolicyName: policyName,
		ExecutionID:       executionID,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
