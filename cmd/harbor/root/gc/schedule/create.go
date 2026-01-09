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
package gcschedule

import (
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/harbor-cli/pkg/api"
	view "github.com/goharbor/harbor-cli/pkg/views/gc/create"
	"github.com/spf13/cobra"
)

func CreateGCScheduleCmd() *cobra.Command {
	var createView view.CreateView
	var nextScheduledTime string
	var params []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create the GC schedule",
		Long: `Create GC schedule.

Examples:
  harbor-cli gc create 
`,

		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			// Validating NextScheduledTime
			if nextScheduledTime != "" {
				dt, err := strfmt.ParseDateTime(nextScheduledTime)
				if err != nil {
					return err
				}

				createView.NextScheduledTime = dt
			}

			// Validating Parameters
			for _, v := range params {
				split := strings.Split(v, "=")

				if len(split) == 2 {
					createView.Parameters[split[0]] = split[1]
				} else {
					return fmt.Errorf("parameter should be of format key=val")
				}
			}

			err = view.CreateScheduleView(&createView)
			if err != nil {
				return err
			}

			err = api.CreateGCSchedule(&createView)
			if err != nil {
				return err
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&createView.Type, "type", "t", "", "type of schedule")
	flags.StringVarP(&createView.Cron, "cron", "c", "", "cron string for when schedule type is 'custom'")
	flags.StringVarP(&nextScheduledTime, "next-schedule", "n", "", "next scheduled time. Example: 2026-01-04T10:15:30Z")
	flags.StringArrayVarP(&params, "parameters", "p", []string{}, "schedule paramters in form key=value")

	return cmd
}
