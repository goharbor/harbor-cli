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

package jobservice

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	view "github.com/goharbor/harbor-cli/pkg/views/jobservice"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func WorkerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Manage workers",
	}

	cmd.AddCommand(ListWorkerCommand())

	return cmd
}

func ListWorkerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [pool-id]",
		Short: "List workers of a pool",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var poolID string
			if len(args) > 0 {
				poolID = args[0]
			} else {
				log.Debug("No pool ID provided, switching to interactive selection...")
				var err error
				poolID, err = view.SelectPoolAsync("Select a Worker Pool")
				if err != nil {
					return err
				}
			}

			if poolID == "all" {
				poolID = ""
			}

			log.Debugf("Attempting to list workers for pool: %s", poolID)

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				log.Debug("Attempting to list workers for formatted output...")
				workers, err := api.ListWorkers(poolID)
				if err != nil {
					return fmt.Errorf("failed to list workers: %v", utils.ParseHarborErrorMsg(err))
				}
				log.WithField("output_format", formatFlag).Debug("Output format selected")
				err = utils.PrintFormat(workers, formatFlag)
				if err != nil {
					return err
				}
			} else {
				err := view.ListWorkersAsync(poolID)
				if err != nil {
					return fmt.Errorf("failed to list workers: %w", err)
				}
			}
			return nil
		},
	}
	return cmd
}
