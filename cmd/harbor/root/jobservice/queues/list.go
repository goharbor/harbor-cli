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
package queues

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	jobserviceutils "github.com/goharbor/harbor-cli/pkg/utils/jobservice"
	queuesview "github.com/goharbor/harbor-cli/pkg/views/jobservice/queues"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListCommand lists all job queues
func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all job queues",
		Long:    "Display all job queues with their pending job counts and latency.",
		Example: "harbor jobservice queues list",
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := api.ListJobQueues()
			if err != nil {
				return jobserviceutils.FormatScheduleError("failed to retrieve job queues", err, "read")
			}

			if response == nil || response.Payload == nil || len(response.Payload) == 0 {
				fmt.Println("No job queues found.")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				return utils.PrintFormat(response.Payload, formatFlag)
			}

			queuesview.ListQueues(response.Payload)
			return nil
		},
	}

	return cmd
}
