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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/gc/update"
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
	var cron string
	var deleteUntagged bool
	var dryRun bool
	var interactive bool

	cmd := &cobra.Command{
		Use:   "update-schedule [schedule-type: none|hourly|daily|weekly|custom]",
		Short: "Update automatic GC schedule",
		Long: `Configure or update the automatic Garbage Collection schedule for the registry.

Available schedule types:
  - none:    Disable scheduled Garbage Collection
  - hourly:  Run GC every hour
  - daily:   Run GC once per day
  - weekly:  Run GC once per week
  - custom:  Define a custom schedule using a cron expression

For custom schedules, Harbor requires a 6-field cron expression in the format:
  seconds minutes hours day-of-month month day-of-week

Examples:
  # Disable automatic Garbage Collection
  harbor gc update-schedule none

  # Configure daily Garbage Collection deleting untagged artifacts
  harbor gc update-schedule daily --delete-untagged

  # Configure custom schedule (e.g. daily at 3:00 AM) in dry-run mode
  harbor gc update-schedule custom --cron "0 0 3 * * *" --dry-run`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scheduleType := cases.Title(language.English).String(strings.ToLower(args[0]))

			if !validScheduleTypes[scheduleType] {
				return fmt.Errorf("invalid schedule type: %s. Valid types are: none, hourly, daily, weekly, custom", args[0])
			}

			if interactive {
				promptCron := (scheduleType == "Custom")
				update.EditGCSchedule(&cron, &deleteUntagged, &dryRun, promptCron)
			}

			logrus.Debugf("Updating GC schedule to type: %s", scheduleType)

			if scheduleType == "Custom" {
				if cron == "" && !interactive {
					return fmt.Errorf("cron expression is required for custom schedule type. Use --cron or --interactive")
				}
				if err := validateCron(cron); err != nil {
					return err
				}
			}

			err := api.UpdateGCSchedule(scheduleType, cron, deleteUntagged, dryRun)
			if err != nil {
				return fmt.Errorf("failed to update GC schedule: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Printf("Successfully updated Garbage Collection schedule to %s\n", scheduleType)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&cron, "cron", "", "Cron expression for custom schedule (include in double quotes)")
	flags.BoolVar(&deleteUntagged, "delete-untagged", false, "Delete untagged artifacts")
	flags.BoolVar(&dryRun, "dry-run", false, "Simulate the GC process without deleting actual blobs")
	flags.BoolVarP(&interactive, "interactive", "i", false, "Update GC schedule interactively")

	return cmd
}

func validateCron(cron string) error {
	if cron == "" {
		return errors.New("cron expression cannot be empty")
	}
	fields := strings.Fields(cron)
	if len(fields) < 6 {
		if len(fields) == 5 {
			return fmt.Errorf("harbor requires 6-field cron format (including seconds). Try: '0 %s'", cron)
		}
		return fmt.Errorf("harbor requires 6-field cron format (seconds minute hour day month weekday)")
	}
	if len(fields) > 6 {
		return fmt.Errorf("too many fields in cron expression, expected 6 but got %d", len(fields))
	}
	return nil
}
