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
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func GCLogsCommand() *cobra.Command {
	var jobID int64

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get logs of a GC Job by id or jobName",
		Long: `Get logs of a GC Job by id or jobName 
Displays key metadata including job kind, schedule, creation time and more.

Examples:
  harbor-cli gc logs # Interactive prompt for GC selection
  harbor-cli gc logs abcd 
`,

		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				id  int64
				err error
			)

			if len(args) == 0 || args[0] == "" {
				id, err = prompt.GetGCJobIDFromUser()
				if err != nil {
					return err
				}
			} else if args[0] != "" {
				id, err = api.GetGCIDFromName(args[0])
				if err != nil {
					return err
				}
			} else if jobID != -1 {
				id = jobID
			} else {
				return fmt.Errorf("no valid id/name specifier found for gc job")
			}

			logs, err := api.LogsGC(id)
			if err != nil {
				return err
			}

			fmt.Println(logs.GetPayload())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&jobID, "id", "i", -1, "Job ID")

	return cmd
}
