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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/gc"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListGCHistory(opts ListFlags) ([]*models.GCHistory, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.GC.GetGCHistory(ctx, &gc.GetGCHistoryParams{
		Page:     &opts.Page,
		PageSize: &opts.PageSize,
		Q:        &opts.Q,
		Sort:     &opts.Sort,
	})
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func GetGCStatus(gcID int64) (*models.GCHistory, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.GC.GetGC(ctx, &gc.GetGCParams{
		GCID: gcID,
	})
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func GetGCLog(gcID int64) (string, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return "", err
	}

	response, err := client.GC.GetGCLog(ctx, &gc.GetGCLogParams{
		GCID: gcID,
	})
	if err != nil {
		return "", err
	}

	return response.Payload, nil
}

func StopGC(gcID int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.GC.StopGC(ctx, &gc.StopGCParams{
		GCID: gcID,
	})
	return err
}

func GetGCSchedule() (*models.GCHistory, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.GC.GetGCSchedule(ctx, &gc.GetGCScheduleParams{})
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func UpdateGCSchedule(scheduleType string, cronStr string, deleteUntagged bool, dryRun bool) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	parameters := map[string]interface{}{
		"delete_untagged": deleteUntagged,
		"dry_run":         dryRun,
	}

	scheduleObj := &models.ScheduleObj{
		Type: scheduleType,
	}
	if scheduleType == "Custom" || scheduleType == "Schedule" {
		scheduleObj.Cron = cronStr
	} else {
		// API requires some placeholder values even for Hourly/Daily/Weekly/None
		scheduleObj.Cron = "0 0 * * * * "
	}

	_, err = client.GC.UpdateGCSchedule(ctx, &gc.UpdateGCScheduleParams{
		Schedule: &models.Schedule{
			Parameters: parameters,
			Schedule:   scheduleObj,
		},
	})
	return err
}

func TriggerGC(deleteUntagged bool, dryRun bool) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	parameters := map[string]interface{}{
		"delete_untagged": deleteUntagged,
		"dry_run":         dryRun,
	}

	scheduleObj := &models.ScheduleObj{
		Type: "Manual",
		Cron: "0 0 * * * * ", // API needs cron
	}

	_, err = client.GC.UpdateGCSchedule(ctx, &gc.UpdateGCScheduleParams{
		Schedule: &models.Schedule{
			Parameters: parameters,
			Schedule:   scheduleObj,
		},
	})
	if err != nil {
		// Fallback to Create if Update fails
		_, err = client.GC.CreateGCSchedule(ctx, &gc.CreateGCScheduleParams{
			Schedule: &models.Schedule{
				Parameters: parameters,
				Schedule:   scheduleObj,
			},
		})
	}

	return err
}
