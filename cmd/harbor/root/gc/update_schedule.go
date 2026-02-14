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

package gc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/gc/update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func UpdateGCScheduleCommand() *cobra.Command {
	var scheduleType string
	var cron string

	cmd := &cobra.Command{
		Use:   "update-schedule",
		Short: "update-schedule [schedule-type: none|hourly|daily|weekly|custom]",
		Long: `Configure or update the automatic GC schedule.

Available schedule types:
  - none:    Disable automatic GC
  - hourly:  Run GC every hour
  - daily:   Run GC once per day
  - weekly:  Run GC once per week
  - custom:  Define a custom schedule using a cron expression`,
		Aliases: []string{"us"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scheduleType = cases.Title(language.English).String(strings.ToLower(args[0]))

			logrus.Infof("Updating GC schedule to type: %s", scheduleType)

			switch scheduleType {
			case "None":
				return updateGCScheduleToNone()
			case "Hourly", "Daily", "Weekly":
				return updateGCPredefinedSchedule(scheduleType)
			case "Custom":
				return updateGCCustomSchedule(cron)
			default:
				return fmt.Errorf("invalid schedule type: %s. Valid types are: none, hourly, daily, weekly, custom", args[0])
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&cron, "cron", "", "Cron expression for custom schedule (include the expression in double quotes)")

	return cmd
}

func updateGCScheduleToNone() error {
	schedule := &models.Schedule{
		Schedule: &models.ScheduleObj{Type: "None"},
	}
	err := api.UpdateGCSchedule(schedule)
	if err != nil {
		return fmt.Errorf("failed to disable GC schedule: %v", utils.ParseHarborErrorMsg(err))
	}
	logrus.Info("Successfully disabled GC schedule")
	return nil
}

func updateGCPredefinedSchedule(scheduleType string) error {
	schedule := &models.Schedule{
		Schedule: &models.ScheduleObj{
			Type: scheduleType,
		},
	}

	err := api.UpdateGCSchedule(schedule)
	if err != nil {
		return fmt.Errorf("failed to update GC schedule: %v", utils.ParseHarborErrorMsg(err))
	}
	logrus.Info("Successfully updated GC schedule")
	return nil
}

func updateGCCustomSchedule(cron string) error {
	if cron == "" {
		logrus.Info("Opening interactive form for custom schedule configuration")
		update.UpdateSchedule(&cron)
		// re-validate after interactive input
		var err error
		cron, err = validateCron(cron)
		if err != nil {
			return err
		}
	}

	schedule := &models.Schedule{
		Schedule: &models.ScheduleObj{
			Type: "Custom",
			Cron: cron,
		},
	}

	err := api.UpdateGCSchedule(schedule)
	if err != nil {
		errMsg := utils.ParseHarborErrorMsg(err)
		if strings.Contains(errMsg, "400") {
			return fmt.Errorf("invalid cron expression: Harbor rejected the schedule. Use the standard 6-field format (seconds minute hour day month weekday)")
		}
		return fmt.Errorf("failed to update GC schedule: %v", errMsg)
	}
	logrus.Infof("Successfully updated GC schedule with custom cron expression: %s", cron)
	return nil
}

func validateCron(cron string) (string, error) {
	if cron == "" {
		return "", errors.New("cron expression cannot be empty")
	}
	fields := strings.Fields(cron)
	if len(fields) < 6 {
		if len(fields) == 5 {
			logrus.Infof("Converting 5-field cron to 6-field by adding '0' for seconds")
			return fmt.Sprintf("0 %s", cron), nil
		}
		return "", fmt.Errorf("harbor requires 6-field cron format (seconds minute hour day month weekday)")
	}
	if len(fields) > 6 {
		return "", fmt.Errorf("too many fields in cron expression, expected 6 but got %d", len(fields))
	}
	return cron, nil
}
