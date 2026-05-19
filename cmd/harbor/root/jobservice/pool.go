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

func PoolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool",
		Short: "Manage worker pools",
	}

	cmd.AddCommand(ListPoolCommand())

	return cmd
}

func ListPoolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the worker pools",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				log.Debug("Attempting to list worker pools for formatted output...")
				pools, err := api.ListWorkerPools()
				if err != nil {
					return fmt.Errorf("failed to list worker pools: %v", utils.ParseHarborErrorMsg(err))
				}
				log.WithField("output_format", formatFlag).Debug("Output format selected")
				err = utils.PrintFormat(pools, formatFlag)
				if err != nil {
					return err
				}
			} else {
				err := view.ListWorkerPoolsAsync()
				if err != nil {
					return fmt.Errorf("failed to list worker pools: %w", err)
				}
			}
			return nil
		},
	}
	return cmd
}
