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

import "github.com/spf13/cobra"

// SchedulesCommand creates the schedules subcommand
func SchedulesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedules",
		Short: "Manage schedules (list/status/pause-all/resume-all)",
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
