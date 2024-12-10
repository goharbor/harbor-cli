package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/schedule"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListSchedule(opts ...ListFlags) (schedule.ListSchedulesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return schedule.ListSchedulesOK{}, err
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
		return schedule.ListSchedulesOK{}, err
	}

	return *response, nil

}
