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
package scan_all

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/scan-all/update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var validScheduleTypes = map[string]bool{
	"None":   true,
	"Hourly": true,
	"Daily":  true,
	"Weekly": true,
	"Custom": true,
}

func UpdateScanAllScheduleCommand() *cobra.Command {
	var scheduleType string
	var cron string

	cmd := &cobra.Command{
		Use:   "update-schedule",
		Short: "update-schedule [schedule-type: none|hourly|daily|weekly|custom]",
		Long: `Configure or update the automatic vulnerability scan schedule for all artifacts.

This command allows you to set when Harbor automatically scans all artifacts for vulnerabilities. You can choose from predefined schedules or create a custom schedule using cron expressions.

Available schedule types:
  - none:    Disable automatic scanning
  - hourly:  Run scan every hour
  - daily:   Run scan once per day
  - weekly:  Run scan once per week
  - custom:  Define a custom schedule using a cron expression

For custom schedules, Harbor requires a 6-field cron expression in the format:
  seconds minutes hours day-of-month month day-of-week

Examples:
  # Disable scheduled scanning
  harbor-cli scan-all update-schedule none

  # Set daily automatic scanning
  harbor-cli scan-all update-schedule daily

  # Set weekly automatic scanning
  harbor-cli scan-all update-schedule weekly

  # Set a custom schedule (every day at 2:30 AM)
  harbor-cli scan-all update-schedule custom --cron "0 30 2 * * *"

  # Use interactive mode to configure a custom schedule
  harbor-cli scan-all update-schedule custom

Note: For custom schedules, if you provide a 5-field cron expression, the CLI will automatically add a leading "0" for the seconds field to create the required 6-field format.`,
		Aliases: []string{"us"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scheduleType = cases.Title(language.English).String(strings.ToLower(args[0]))

			if !validScheduleTypes[scheduleType] {
				return fmt.Errorf("invalid schedule type: %s. Valid types are: none, hourly, daily, weekly, custom", args[0])
			}

			logrus.Infof("Updating scan all schedule to type: %s", scheduleType)

			switch scheduleType {
			case "None":
				return updateScheduleToNone()

			case "Hourly", "Daily", "Weekly":
				return updatePredefinedSchedule(scheduleType)

			case "Custom":
				return updateCustomSchedule(cron)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&cron, "cron", "", "Cron expression for custom schedule (include the expression in double quotes)")

	return cmd
}

func updateScheduleToNone() error {
	logrus.Info("Setting scan all schedule to None (disabled)")
	err := api.UpdateScanAllSchedule(models.ScheduleObj{Type: "None"})
	if err != nil {
		return fmt.Errorf("failed to disable scan schedule: %v", utils.ParseHarborErrorMsg(err))
	}
	logrus.Info("Successfully disabled scan all schedule")
	return nil
}

func updatePredefinedSchedule(scheduleType string) error {
	logrus.Infof("Setting scan all schedule to %s", scheduleType)

	// Random cron expression and time needed by API
	randomCron := "0 0 * * * * "
	randomTime := strfmt.DateTime{}

	err := api.UpdateScanAllSchedule(models.ScheduleObj{
		Type:              scheduleType,
		Cron:              randomCron,
		NextScheduledTime: randomTime,
	})

	if err != nil {
		return fmt.Errorf("failed to update scan schedule: %v", utils.ParseHarborErrorMsg(err))
	}

	logrus.Infof("Successfully set scan all schedule to %s", scheduleType)
	return nil
}

func updateCustomSchedule(cron string) error {
	if cron == "" {
		logrus.Info("Opening interactive form for custom schedule configuration")
		update.UpdateSchedule(&cron)
	}

	if err := validateCron(cron); err != nil {
		return err
	}

	logrus.Infof("Setting scan all schedule with custom cron expression: %s", cron)

	// Random time needed by API
	randomTime := strfmt.DateTime{}
	err := api.UpdateScanAllSchedule(models.ScheduleObj{
		Type:              "Custom",
		Cron:              cron,
		NextScheduledTime: randomTime,
	})

	if err != nil {
		errMsg := utils.ParseHarborErrorMsg(err)
		if strings.Contains(errMsg, "400") {
			return fmt.Errorf("invalid cron expression: Harbor rejected the schedule. Use the standard 5-field format (minute hour day month weekday)")
		}
		return fmt.Errorf("failed to update scan schedule: %v", errMsg)
	}

	logrus.Info("Successfully set scan all schedule with custom cron expression")
	return nil
}

func validateCron(cron string) error {
	if cron == "" {
		return errors.New("cron expression cannot be empty")
	}
	fields := strings.Fields(cron)
	if len(fields) < 6 {
		if len(fields) == 5 {
			logrus.Infof("Converting 5-field cron to 6-field by adding '0' for seconds")
			return fmt.Errorf("harbor requires 6-field cron format (including seconds). Try: '0 %s'", cron)
		}
		return fmt.Errorf("harbor requires 6-field cron format (seconds minute hour day month weekday)")
	}
	if len(fields) > 6 {
		return fmt.Errorf("too many fields in cron expression, expected 6 but got %d", len(fields))
	}
	return nil
}
