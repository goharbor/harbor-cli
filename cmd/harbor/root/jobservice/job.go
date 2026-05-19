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

package jobservice

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/logviewer"
	view "github.com/goharbor/harbor-cli/pkg/views/jobservice"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func JobCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "Manage individual jobs",
	}

	cmd.AddCommand(
		StopJobCommand(),
		LogJobCommand(),
	)

	return cmd
}

func StopJobCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop [job-id]",
		Short: "Stop a particular job",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var jobID string
			if len(args) > 0 {
				jobID = args[0]
			} else {
				log.Debug("No job ID provided, switching to interactive selection...")
				var err error
				jobID, err = view.SelectRunningJobAsync()
				if err != nil {
					return err
				}
			}

			log.Debugf("Attempting to stop job: %s", jobID)
			err := api.StopJob(jobID)
			if err != nil {
				return fmt.Errorf("failed to stop job: %v", utils.ParseHarborErrorMsg(err))
			}
			fmt.Printf("Job \"%s\" stopped successfully\n", jobID)
			return nil
		},
	}
	return cmd
}

func LogJobCommand() *cobra.Command {
	var follow bool
	var refreshInterval string

	cmd := &cobra.Command{
		Use:   "log <job-id>",
		Short: "Display logs of a particular job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jobID := args[0]

			interval := 5 * time.Second
			if refreshInterval != "" {
				var err error
				interval, err = time.ParseDuration(refreshInterval)
				if err != nil {
					return fmt.Errorf("invalid refresh interval: %w", err)
				}
			}

			m := logviewer.NewModel(jobID, api.GetJobLog, follow, interval)
			if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
				return fmt.Errorf("error running log viewer: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	cmd.Flags().StringVarP(&refreshInterval, "refresh-interval", "n", "", "Interval to refresh logs (default 5s)")

	return cmd
}
