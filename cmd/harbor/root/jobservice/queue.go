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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	view "github.com/goharbor/harbor-cli/pkg/views/jobservice"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func QueueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queue",
		Short: "Manage job queues",
	}

	cmd.AddCommand(
		ListQueueCommand(),
		ClearQueueCommand(),
		PauseQueueCommand(),
		ResumeQueueCommand(),
	)

	return cmd
}

func ListQueueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List job queues",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				log.Debug("Attempting to list job queues for formatted output...")
				queues, err := api.ListJobQueues()
				if err != nil {
					return fmt.Errorf("failed to list job queues: %v", utils.ParseHarborErrorMsg(err))
				}
				log.WithField("output_format", formatFlag).Debug("Output format selected")
				err = utils.PrintFormat(queues, formatFlag)
				if err != nil {
					return err
				}
			} else {
				err := view.ListJobQueuesAsync()
				if err != nil {
					return fmt.Errorf("failed to list job queues: %w", err)
				}
			}
			return nil
		},
	}
	return cmd
}

func ClearQueueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear [job-type]",
		Short: "Clear a particular job queue",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var jobType string
			if len(args) > 0 {
				jobType = args[0]
			} else {
				log.Debug("No job type provided, switching to interactive selection...")
				var err error
				jobType, err = view.SelectQueueAsync("Select a Queue to Clear")
				if err != nil {
					return err
				}
			}

			log.Debugf("Attempting to clear job queue: %s", jobType)
			err := api.ActionPendingJobs(jobType, api.JobActionStop)
			if err != nil {
				return fmt.Errorf("failed to clear job queue: %v", utils.ParseHarborErrorMsg(err))
			}
			fmt.Printf("Pending jobs in jobservice queue \"%s\" cleared successfully\n", jobType)
			return nil
		},
	}
	return cmd
}

func PauseQueueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause [job-type]",
		Short: "Pause a particular job queue",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var jobType string
			if len(args) > 0 {
				jobType = args[0]
			} else {
				log.Debug("No job type provided, switching to interactive selection...")
				var err error
				jobType, err = view.SelectQueueAsync("Select a Queue to Pause")
				if err != nil {
					return err
				}
			}

			log.Debugf("Attempting to pause job queue: %s", jobType)
			err := api.ActionPendingJobs(jobType, api.JobActionPause)
			if err != nil {
				return fmt.Errorf("failed to pause job queue: %v", utils.ParseHarborErrorMsg(err))
			}
			fmt.Printf("Jobservice queue \"%s\" paused successfully\n", jobType)
			return nil
		},
	}
	return cmd
}

func ResumeQueueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume [job-type]",
		Short: "Resume a particular job queue",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var jobType string
			if len(args) > 0 {
				jobType = args[0]
			} else {
				log.Debug("No job type provided, switching to interactive selection...")
				var err error
				jobType, err = view.SelectQueueAsync("Select a Queue to Resume")
				if err != nil {
					return err
				}
			}

			log.Debugf("Attempting to resume job queue: %s", jobType)
			err := api.ActionPendingJobs(jobType, api.JobActionResume)
			if err != nil {
				return fmt.Errorf("failed to resume job queue: %v", utils.ParseHarborErrorMsg(err))
			}
			fmt.Printf("Jobservice queue \"%s\" resumed successfully\n", jobType)
			return nil
		},
	}
	return cmd
}
