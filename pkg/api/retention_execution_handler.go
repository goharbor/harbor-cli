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
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/retention"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListRetentionExecutions(retentionID string) (*retention.ListRetentionExecutionsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	policyID, err := strconv.ParseInt(retentionID, 10, 64)
	if err != nil {
		return nil, err
	}

	return client.Retention.ListRetentionExecutions(ctx, &retention.ListRetentionExecutionsParams{ID: policyID})
}

func ListRetentionTasks(retentionID string, executionID int64) (*retention.ListRetentionTasksOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	policyID, err := strconv.ParseInt(retentionID, 10, 64)
	if err != nil {
		return nil, err
	}

	return client.Retention.ListRetentionTasks(ctx, &retention.ListRetentionTasksParams{ID: policyID, Eid: executionID})
}
