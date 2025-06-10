package scan_all

import (
	"errors"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/scan-all/update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
				logrus.Infof("Updating scan all schedule to type: %s", scheduleType)

				if scheduleType == "None" {
					logrus.Info("Setting scan all schedule to None (disabled)")
					err := api.UpdateScanAllSchedule(models.ScheduleObj{Type: "None"})
					if err != nil {
						logrus.Errorf("Failed to update scan schedule: %v", err)
						return err
					}
					logrus.Info("Successfully disabled scan all schedule")
					return nil
				} else if scheduleType == "Hourly" || scheduleType == "Daily" || scheduleType == "Weekly" {
					logrus.Infof("Setting scan all schedule to %s", scheduleType)
					// Random cron expression and random time need to be passed to the API, even though they are not used, otherwise it returns bad request
					randomCron := "0 * * * * *"
					randomTime := strfmt.DateTime{}
					err := api.UpdateScanAllSchedule(models.ScheduleObj{Type: scheduleType, Cron: randomCron, NextScheduledTime: randomTime})
					if err != nil {
						logrus.Errorf("Failed to update scan schedule: %v", err)
						return err
					}
					logrus.Infof("Successfully set scan all schedule to %s", scheduleType)
					return nil
				} else if scheduleType == "Custom" {
					if cron != "" {
						logrus.Infof("Setting scan all schedule to Custom with cron expression: %s", cron)
						// Random time need to be passed to the API, same reason as above
						randomTime := strfmt.DateTime{}
						err := api.UpdateScanAllSchedule(models.ScheduleObj{Type: "Schedule", Cron: cron, NextScheduledTime: randomTime})
						if err != nil {
							logrus.Errorf("Failed to update scan schedule: %v", err)
							return err
						}
						logrus.Info("Successfully set scan all schedule with custom cron expression")
						return nil
					} else {
						logrus.Info("Opening interactive form for custom schedule configuration")
						update.UpdateSchedule(&cron)
						logrus.Infof("Setting scan all schedule with custom cron expression: %s", cron)
						// Random time need to be passed to the API, same reason as above
						randomTime := strfmt.DateTime{}
						err := api.UpdateScanAllSchedule(models.ScheduleObj{Type: "Schedule", Cron: cron, NextScheduledTime: randomTime})
						if err != nil {
							logrus.Errorf("Failed to update scan schedule: %v", err)
							return err
						}
						logrus.Info("Successfully set scan all schedule with custom cron expression")
						return nil
					}
				} else {
					logrus.Errorf("Invalid schedule type: %s", scheduleType)
					return errors.New("invalid schedule type")
				}
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&cron, "cron", "", "Cron expression (include the expression in double quotes)")

	return cmd
}
