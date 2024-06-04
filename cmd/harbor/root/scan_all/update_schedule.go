package scan_all

import (
	"errors"
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/scan-all/update"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func UpdateScanAllScheduleCommand() *cobra.Command {
	var scheduleType string
	var cron string

	cmd := &cobra.Command{
		Use:     "update-schedule",
		Short:   "update-schedule [schedule-type: none|hourly|daily|weekly|custom]",
		Aliases: []string{"us"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("schedule type is required")
			} else if len(args) > 1 {
				return errors.New("too many arguments")
			} else {
				scheduleType = cases.Title(language.English).String(strings.ToLower(args[0]))
				if scheduleType == "None" {
					return api.UpdateScanAllSchedule(models.ScheduleObj{Type: "None"})
				} else if scheduleType == "Hourly" || scheduleType == "Daily" || scheduleType == "Weekly" {
					// Random cron expression and random time need to be passed to the API, even though they are not used, otherwise it returns bad request
					randomCron := "0 * * * * *"
					randomTime := strfmt.DateTime{}
					return api.UpdateScanAllSchedule(models.ScheduleObj{Type: scheduleType, Cron: randomCron, NextScheduledTime: randomTime})
				} else if scheduleType == "Custom" {
					if cron != "" {
						// Random time need to be passed to the API, same reason as above
						randomTime := strfmt.DateTime{}
						return api.UpdateScanAllSchedule(models.ScheduleObj{Type: "Schedule", Cron: cron, NextScheduledTime: randomTime})
					} else {
						update.UpdateSchedule(&cron)
						// Random time need to be passed to the API, same reason as above
						randomTime := strfmt.DateTime{}
						return api.UpdateScanAllSchedule(models.ScheduleObj{Type: "Schedule", Cron: cron, NextScheduledTime: randomTime})
					}
				} else {
					return errors.New("invalid schedule type")
				}
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&cron, "cron", "", "Cron expression (include the expression in double quotes)")

	return cmd
}
