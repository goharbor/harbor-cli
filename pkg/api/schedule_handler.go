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
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/schedule"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListSchedule(opts ...ListFlags) (schedule.ListSchedulesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return schedule.ListSchedulesOK{}, fmt.Errorf("failed to initialize client context for schedules")
	}

	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Schedule.ListSchedules(ctx, &schedule.ListSchedulesParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
	})

	if err != nil {
		switch err.(type) {
		case *schedule.ListSchedulesInternalServerError:
			return schedule.ListSchedulesOK{}, fmt.Errorf("internal server error occurred while listing schedules")
		default:
			return schedule.ListSchedulesOK{}, fmt.Errorf("unknown error occurred while listing schedules")
		}
	}

	return *response, nil
}
