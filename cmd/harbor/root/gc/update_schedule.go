package gc

import (
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
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

			if !validScheduleTypes[scheduleType] {
				return fmt.Errorf("invalid schedule type: %s. Valid types are: none, hourly, daily, weekly, custom", args[0])
			}

			logrus.Infof("Updating GC schedule to type: %s", scheduleType)

			switch scheduleType {
			case "None":
				return updateGCScheduleToNone()
			case "Hourly", "Daily", "Weekly":
				return updateGCPredefinedSchedule(scheduleType)
			case "Custom":
				return updateGCCustomSchedule(cron)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&cron, "cron", "", "Cron expression for custom schedule (include the expression in double quotes)")

	return cmd
}

func updateGCScheduleToNone() error {
	// Wrap ScheduleObj in proper structure
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
	randomCron := "0 0 * * * * "
	randomTime := strfmt.DateTime{}

	schedule := &models.Schedule{
		Schedule: &models.ScheduleObj{
			Type:              scheduleType,
			Cron:              randomCron,
			NextScheduledTime: randomTime,
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
		return fmt.Errorf("cron expression is required for custom schedule")
	}

	randomTime := strfmt.DateTime{}

	schedule := &models.Schedule{
		Schedule: &models.ScheduleObj{
			Type:              "Custom",
			Cron:              cron,
			NextScheduledTime: randomTime,
		},
	}

	err := api.UpdateGCSchedule(schedule)
	if err != nil {
		return fmt.Errorf("failed to update GC schedule: %v", utils.ParseHarborErrorMsg(err))
	}
	logrus.Info("Successfully updated GC schedule with custom cron")
	return nil
}
