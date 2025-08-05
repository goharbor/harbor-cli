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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func LogsCommand() *cobra.Command {
	var taskID int64
	var execID int64
	cmd := &cobra.Command{
		Use:   "log [EXECUTION_ID] [TASK_ID]",
		Short: "get replication execution logs by execution and task id",
		Long:  `Get the logs of a specific replication execution and task by their IDs. If no IDs are provided, it will prompt the user to select them interactively.`,
		Example: `  harbor replication log -e 12345 -t 67890
		  harbor replication log --execution-id 12345 --task-id 67890
		  harbor replication log --execution-id 12345
  harbor replication log`,
		Args: cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if execID != 0 && taskID == 0 {
				taskID = prompt.GetReplicationTaskIDFromUser(execID)
			} else if execID == 0 && taskID == 0 {
				rpolicyID := prompt.GetReplicationPolicyFromUser()
				execID = prompt.GetReplicationExecutionIDFromUser(rpolicyID)
				taskID = prompt.GetReplicationTaskIDFromUser(execID)
			} else if execID == 0 && taskID != 0 {
				return fmt.Errorf("execution ID must be provided if task ID is specified")
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

	flags := cmd.Flags()
	flags.Int64VarP(&execID, "execution-id", "e", 0, "Replication execution ID")
	flags.Int64VarP(&taskID, "task-id", "t", 0, "Replication task ID")

	return cmd
}
