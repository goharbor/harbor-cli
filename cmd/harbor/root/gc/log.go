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
	gclog "github.com/goharbor/harbor-cli/pkg/views/gc/log"
	"github.com/spf13/cobra"
)

func GetGCLogCommand() *cobra.Command {
	var gcID int64

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Get GC job log",
		Long: `Get the log of a specific GC (Garbage Collection) job.

If no GC job ID is provided via the --id flag, an interactive selector
will be displayed to choose from available GC jobs.

Examples:
  # Get GC log by specifying the job ID
  harbor gc log --id 42

  # Get GC log interactively (select from list)
  harbor gc log`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if gcID <= 0 {
				gcID, err = gclog.SelectGCJob()
				if err != nil {
					return err
				}
			}

			logData, err := api.GetGCJobLog(gcID)
			if err != nil {
				return fmt.Errorf("failed to get GC log: %v", err)
			}

			fmt.Println(logData)
			return nil
		},
	}

	cmd.Flags().Int64Var(&gcID, "id", 0, "ID of the GC job to get logs for")

	return cmd
}
