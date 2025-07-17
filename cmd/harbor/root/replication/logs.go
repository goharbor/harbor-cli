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
package replication

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func LogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log [EXECUTION_ID] [TASK_ID]",
		Short: "get replication execution logs by execution and task id",
		Long:  `Get the logs of a specific replication execution and task by their IDs. If no IDs are provided, it will prompt the user to select them interactively.`,
		Example: `  harbor replication log 12345 67890
  harbor replication log`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var taskID int64
			var execID int64
			if len(args) == 2 {
				var err error
				execID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid replication execution ID: %s, %v", args[0], err)
				}
				taskID, err = strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid replication task ID: %s, %v", args[1], err)
				}
			} else {
				rpolicyID := prompt.GetReplicationPolicyFromUser()
				execID = prompt.GetReplicationExecutionIDFromUser(rpolicyID)
				taskID = prompt.GetReplicationTaskIDFromUser(execID)
			}

			logs, err := api.GetReplicationLog(execID, taskID)
			if err != nil {
				return fmt.Errorf("failed to get replication task logs: %v", utils.ParseHarborErrorMsg(err))
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(logs.Payload, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				if logs.Payload != "" {
					fmt.Println(logs.Payload)
				}
			}
			return nil
		},
	}
	return cmd
}
