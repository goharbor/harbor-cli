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
	"github.com/spf13/cobra"
)

func GetGCLogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log [gc_id]",
		Short: "Get GC job log",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("Invalid GC ID: %v", err)
			}

			if id <= 0 {
				return fmt.Errorf("gc_id must be a positive integer")
			}

			logData, err := api.GetGCJobLog(id)
			if err != nil {
				return fmt.Errorf("Failed to get GC log: %v", err)
			}

			fmt.Println(logData)
			return nil
		},
	}
	return cmd
}
