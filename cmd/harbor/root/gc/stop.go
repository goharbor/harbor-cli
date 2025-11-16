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
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func GCStopCommand() *cobra.Command {
	var jobID int64

	cmd := &cobra.Command{
		Use:   "history",
		Short: "List cleanup jobs in a Harbor instance",
		Long: `List all cleanup jobs in a Harbor instance. 
Displays key metadata including job kind, schedule, creation time and more.

Examples:
  harbor-cli cleanup list                # Interactive prompt for project and repository

Supports pagination, search queries, and sorting using flags.`,

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
				// TODO: Add Name to ID Conversion

				strID, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}

				id = int64(strID)
			} else if jobID != -1 {
				id = jobID
			} else {
				return fmt.Errorf("no valid id/name specifier found for gc job")
			}

			err = api.StopGC(id)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully stopped job with ID: %d \n", id)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&jobID, "id", "i", -1, "Job ID")

	return cmd
}
