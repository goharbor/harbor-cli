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
	log "github.com/sirupsen/logrus"
)

// GetGCHistory gets the GC history
func GetGCHistory(opts ListFlags) ([]*models.GCHistory, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	params := &gc.GetGCHistoryParams{
		Page:     &opts.Page,
		PageSize: &opts.PageSize,
		Q:        &opts.Q,
		Sort:     &opts.Sort,
	}

	resp, err := client.GC.GetGCHistory(ctx, params)
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

// GetGCJobLog gets the log of a specific GC job
func GetGCJobLog(id int64) (string, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return "", err
	}

	resp, err := client.GC.GetGCLog(ctx, &gc.GetGCLogParams{
		GCID: id,
	})
	if err != nil {
		return "", err
	}

	return resp.Payload, nil
}

// GetGCSchedule gets the GC schedule
func GetGCSchedule() (*models.GCHistory, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.GC.GetGCSchedule(ctx, &gc.GetGCScheduleParams{})
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

// CreateGCSchedule creates a GC schedule
func CreateGCSchedule(schedule *models.Schedule) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.GC.CreateGCSchedule(ctx, &gc.CreateGCScheduleParams{
		Schedule: schedule,
	})

	if err != nil {
		return err
	}

	log.Info("GC schedule created successfully")
	return nil
}

// UpdateGCSchedule updates the GC schedule
// Modified to take *models.Schedule allowing passing parameters
func UpdateGCSchedule(schedule *models.Schedule) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.GC.UpdateGCSchedule(ctx, &gc.UpdateGCScheduleParams{
		Schedule: schedule,
	})

	if err != nil {
		return err
	}

	log.Info("GC schedule updated successfully")
	return nil
}

// StopGC stops a running GC job by ID
func StopGC(gcID int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.GC.StopGC(ctx, &gc.StopGCParams{
		GCID: gcID,
	})

	if err != nil {
		return err
	}

	log.Infof("GC job %d stopped successfully", gcID)
	return nil
}
