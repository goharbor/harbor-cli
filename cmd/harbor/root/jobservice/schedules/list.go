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
package schedules

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/jobservice/schedules"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SchedulesCommand creates the schedules subcommand
func SchedulesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedules",
		Short: "Manage job schedules",
		Long:  "List schedules and manage global scheduler status.",
	}

	cmd.AddCommand(
		ListCommand(),
		StatusCommand(),
		PauseAllCommand(),
		ResumeAllCommand(),
	)

	return cmd
}

// ListCommand lists all schedules
func ListCommand() *cobra.Command {
	var page int64 = 1
	var pageSize int64 = 20

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all schedules",
		Long:    "Display all job schedules with pagination support.",
		Example: "harbor jobservice schedules list --page 1 --page-size 20",
		RunE: func(cmd *cobra.Command, args []string) error {
			if page < 1 {
				return fmt.Errorf("page must be >= 1")
			}
			if pageSize < 1 || pageSize > 100 {
				return fmt.Errorf("page-size must be between 1 and 100")
			}

			response, err := api.ListSchedules(page, pageSize)
			if err != nil {
				return fmt.Errorf("failed to retrieve schedules: %w", err)
			}

			if response == nil || response.Payload == nil || len(response.Payload) == 0 {
				fmt.Println("No schedules found.")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				return utils.PrintFormat(response.Payload, formatFlag)
			}

			totalCount := response.XTotalCount
			schedules.ListSchedules(response.Payload, page, pageSize, totalCount)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&page, "page", 1, "Page number")
	flags.Int64Var(&pageSize, "page-size", 20, "Number of items per page")

	return cmd
}

// StatusCommand shows the global scheduler status
func StatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Show scheduler status",
		Long:    "Display whether the global scheduler is paused or running.",
		Example: "harbor jobservice schedules status",
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := api.GetSchedulePaused()
			if err != nil {
				return fmt.Errorf("failed to retrieve scheduler status: %w", err)
			}

			if response == nil || response.Payload == nil {
				fmt.Println("Unable to determine scheduler status.")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				return utils.PrintFormat(response.Payload, formatFlag)
			}

			schedules.PrintScheduleStatus(response.Payload)
			return nil
		},
	}

	return cmd
}

// PauseAllCommand pauses all schedules
func PauseAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause-all",
		Short:   "Pause all schedules",
		Long:    "Pause the global scheduler and all schedules.",
		Example: "harbor jobservice schedules pause-all",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Pausing all schedules...")
			err := api.ActionJobQueue("SCHEDULER", "pause")
			if err != nil {
				return fmt.Errorf("failed to pause all schedules: %w", err)
			}
			fmt.Println("✓ All schedules paused successfully.")
			return nil
		},
	}

	return cmd
}

// ResumeAllCommand resumes all schedules
func ResumeAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resume-all",
		Short:   "Resume all schedules",
		Long:    "Resume the global scheduler and all schedules.",
		Example: "harbor jobservice schedules resume-all",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Resuming all schedules...")
			err := api.ActionJobQueue("SCHEDULER", "resume")
			if err != nil {
				return fmt.Errorf("failed to resume all schedules: %w", err)
			}
			fmt.Println("✓ All schedules resumed successfully.")
			return nil
		},
	}

	return cmd
}
