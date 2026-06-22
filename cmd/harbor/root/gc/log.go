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
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LogGCOperation() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "log [gc-id]",
		Short:   "Get GC execution log",
		Long:    `Retrieve the execution log of a specific Garbage Collection run by its ID.`,
		Example: `  harbor gc log 12`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gcID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid GC execution ID: %v", args[0])
			}

			logrus.Debugf("Fetching Garbage Collection execution log for ID: %d", gcID)
			logs, err := api.GetGCLog(gcID)
			if err != nil {
				return fmt.Errorf("failed to get GC log: %v", utils.ParseHarborErrorMsg(err))
			}

			if logs == "" {
				fmt.Printf("No logs found for Garbage Collection run ID: %d\n", gcID)
				return nil
			}

			fmt.Println(logs)
			return nil
		},
	}

	return cmd
}
