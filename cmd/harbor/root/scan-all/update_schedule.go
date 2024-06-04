package scan_all

import (
	"errors"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/scan-all/update"
	"github.com/spf13/cobra"
)

func UpdateScanAllScheduleCommand() *cobra.Command {
	var scheduleType string
	var cron string

	cmd := &cobra.Command{
		Use:     "update-schedule",
		Short:   "update-schedule [schedule-type: None|Hourly|Daily|Weekly|Custom]",
		Aliases: []string{"us"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("schedule type is required")
			} else if len(args) > 1 {
				return errors.New("too many arguments")
			} else {
				scheduleType = args[0]
				if scheduleType == "None" {
					return api.UpdateScanAllSchedule(models.ScheduleObj{Type: "None"})
				} else if scheduleType == "Hourly" || scheduleType == "Daily" || scheduleType == "Weekly" {
					return api.UpdateScanAllSchedule(models.ScheduleObj{Type: scheduleType})
				} else if scheduleType == "Custom" {
					if cron != "" {
						return api.UpdateScanAllSchedule(models.ScheduleObj{Type: "Schedule", Cron: cron})
					} else {
						update.UpdateSchedule(&cron)
						return api.UpdateScanAllSchedule(models.ScheduleObj{Type: "Schedule", Cron: cron})
					}
				} else {
					return errors.New("invalid schedule type")
				}
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&cron, "cron", "", "Cron expression")

	return cmd
}
