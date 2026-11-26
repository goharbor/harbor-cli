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
package workers

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	jobserviceutils "github.com/goharbor/harbor-cli/pkg/utils/jobservice"
	workersview "github.com/goharbor/harbor-cli/pkg/views/jobservice/workers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListCommand lists workers in a pool
func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [pool-id]",
		Short:   "List workers in a pool",
		Long:    "Display all workers in the specified Harbor jobservice worker pool. Requires system admin privileges.",
		Example: "harbor jobservice workers list <pool-id>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			poolID := args[0]
			if poolID == "" {
				return fmt.Errorf("pool-id is required")
			}

			response, err := api.GetWorkers(poolID)
			if err != nil {
				return jobserviceutils.FormatScheduleError("failed to retrieve workers", err, "read")
			}

			if response == nil || response.Payload == nil || len(response.Payload) == 0 {
				fmt.Println("No workers found.")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				return utils.PrintFormat(response.Payload, formatFlag)
			}

			workersview.ListWorkers(response.Payload)
			return nil
		},
	}

	return cmd
}
