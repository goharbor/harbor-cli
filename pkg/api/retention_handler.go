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
	"errors"
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/retention"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

var ErrNoRetentionPolicy = errors.New("no retention policy configured")

func CreateRetention(policy *models.RetentionPolicy) (string, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return "", err
	}

	if policy == nil {
		return "", fmt.Errorf("retention policy payload cannot be nil")
	}

	response, err := client.Retention.CreateRetention(ctx, &retention.CreateRetentionParams{
		Policy: policy,
	})
	if err != nil {
		return "", err
	}

	return response.Location, nil
}

func UpdateRetention(retentionID int64, policy *models.RetentionPolicy) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	if policy == nil {
		return fmt.Errorf("retention policy payload cannot be nil")
	}

	_, err = client.Retention.UpdateRetention(ctx, &retention.UpdateRetentionParams{
		ID:     retentionID,
		Policy: policy,
	})
	return err
}

func TriggerRetentionExecution(retentionID int64, dryRun bool) (string, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return "", err
	}

	_, response, err := client.Retention.TriggerRetentionExecution(ctx, &retention.TriggerRetentionExecutionParams{
		ID: retentionID,
		Body: retention.TriggerRetentionExecutionBody{
			DryRun: dryRun,
		},
	})
	if err != nil {
		return "", err
	}

	if response != nil {
		return response.Location, nil
	}

	return "", nil
}

// GetRetentionPolicy retrieves a retention policy by ID.
func GetRetentionPolicy(retentionID int64) (*models.RetentionPolicy, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	params := retention.NewGetRetentionParams()
	params.SetID(retentionID)

	resp, err := client.Retention.GetRetention(ctx, params)
	if err != nil {
		return nil, err
	}

	return resp.GetPayload(), nil
}

func GetRetentionIDByProjectName(projectName string) (int64, error) {
	if projectName == "" {
		return 0, fmt.Errorf("project name is required")
	}

	response, err := ListConfig(false, projectName)
	if err != nil {
		return 0, err
	}

	if response == nil || response.Payload == nil {
		return 0, fmt.Errorf("%w for project %q", ErrNoRetentionPolicy, projectName)
	}

	retentionIDValue, ok := response.Payload["retention_id"]
	if !ok || retentionIDValue == "" {
		return 0, fmt.Errorf("%w for project %q", ErrNoRetentionPolicy, projectName)
	}

	retentionID, err := strconv.ParseInt(retentionIDValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid retention_id %q for project %q: %w", retentionIDValue, projectName, err)
	}

	return retentionID, nil
}
