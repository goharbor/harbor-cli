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
	jobserviceutils "github.com/goharbor/harbor-cli/pkg/utils/jobservice"
	"github.com/spf13/cobra"
)

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
				return jobserviceutils.FormatScheduleError("failed to resume all schedules", err, "ActionStop")
			}
			fmt.Println("✓ All schedules resumed successfully.")
			return nil
		},
	}

	return cmd
}
