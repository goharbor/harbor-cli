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
package executions

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/replication/execution/view"
	tasklist "github.com/goharbor/harbor-cli/pkg/views/replication/task/list"
)

type ExecutionsAndTasks struct {
	Execution *models.ReplicationExecution      `json:"execution,omitempty" yaml:"execution,omitempty"`
	Tasks     map[int64]*models.ReplicationTask `json:"tasks,omitempty" yaml:"tasks,omitempty"`
}

func NewExecutionsAndTasks(execution *models.ReplicationExecution, tasks []*models.ReplicationTask) *ExecutionsAndTasks {
	taskMap := make(map[int64]*models.ReplicationTask)
	for _, task := range tasks {
		if task.ID != 0 {
			taskMap[task.ID] = task
		}
	}
	return &ExecutionsAndTasks{
		Execution: execution,
		Tasks:     taskMap,
	}
}

func ViewCommand() *cobra.Command {
	var notListTasks bool
	cmd := &cobra.Command{
		Use:   "view [ID]",
		Short: "get replication execution by id",
		Long:  `Get a specific replication execution by its ID. If no ID is provided, it will prompt the user to select one interactively. If the --no-tasks flag is set, it will not list associated tasks.`,
		Example: `  harbor replication executions view 12345
  harbor replication executions view --no-tasks`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var execID int64
			if len(args) > 0 {
				var err error
				// convert string to int64
				execID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid replication execution ID: %s, %v", args[0], err)
				}
			} else {
				rpolicyID := prompt.GetReplicationPolicyFromUser()
				execID = prompt.GetReplicationExecutionIDFromUser(rpolicyID)
			}

			execution, err := api.GetReplicationExecution(execID)
			if err != nil {
				return fmt.Errorf("failed to get replication execution: %v", utils.ParseHarborErrorMsg(err))
			}

			tasks, err := api.ListReplicationTasks(execID)
			if err != nil {
				return fmt.Errorf("failed to list replication tasks: %v", utils.ParseHarborErrorMsg(err))
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				if notListTasks {
					err = utils.PrintFormat(execution.Payload, FormatFlag)
					if err != nil {
						return err
					}
				} else {
					combinedOutput := NewExecutionsAndTasks(execution.Payload, tasks.Payload)
					err = utils.PrintFormat(combinedOutput, FormatFlag)
					if err != nil {
						return err
					}
				}
			} else {
				view.ViewExecution(execution.Payload)
				if !notListTasks {
					fmt.Println("Tasks:")
					tasklist.ListTasks(tasks.Payload)
				}
			}
			return nil
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&notListTasks, "no-tasks", "", false, "Do not list associated tasks")
	return cmd
}
