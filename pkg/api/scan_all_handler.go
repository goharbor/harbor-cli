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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/scan_all"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func CreateScanAllSchedule(schedule models.ScheduleObj) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.ScanAll.CreateScanAllSchedule(ctx, &scan_all.CreateScanAllScheduleParams{Schedule: &models.Schedule{Schedule: &schedule}})

	if err != nil {
		return err
	}

	if response != nil {
		// The CreateScanAllSchedule API is used only for scanning all artifacts now
		log.Info("Scan started successfully")
	}
	return nil
}

func UpdateScanAllSchedule(schedule models.ScheduleObj) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.ScanAll.UpdateScanAllSchedule(ctx, &scan_all.UpdateScanAllScheduleParams{Schedule: &models.Schedule{Schedule: &schedule}})

	if err != nil {
		return err
	}

	if response != nil {
		log.Info("Schedule updated successfully")
	}
	return nil
}

func StopScanAll() error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.ScanAll.StopScanAll(ctx, &scan_all.StopScanAllParams{})

	if err != nil {
		return err
	}

	if response != nil {
		log.Info("Scan all stopped successfully")
	}
	return nil
}

func GetScanAllSchedule() (*models.ScheduleObj, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.ScanAll.GetScanAllSchedule(ctx, &scan_all.GetScanAllScheduleParams{})

	if err != nil {
		return nil, err
	}

	return response.Payload.Schedule, nil
}

func GetScanAllMetrics(scheduled bool) (*models.Stats, error) {
	ctx, client, clientErr := utils.ContextWithClient()
	if clientErr != nil {
		return nil, clientErr
	}

	if scheduled {
		response, responseErr := client.ScanAll.GetLatestScheduledScanAllMetrics(ctx, &scan_all.GetLatestScheduledScanAllMetricsParams{})
		if responseErr != nil {
			return nil, responseErr
		}
		return response.Payload, nil
	} else {
		response, responseErr := client.ScanAll.GetLatestScanAllMetrics(ctx, &scan_all.GetLatestScanAllMetricsParams{})
		if responseErr != nil {
			return nil, responseErr
		}
		return response.Payload, nil
	}
}
