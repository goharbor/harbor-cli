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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func StartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start replication",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("Starting replication")

			var rpolicyID int64
			if len(args) > 0 {
				var err error
				// convert string to int64
				rpolicyID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid replication policy ID: %s, %v", args[0], err)
				}
			} else {
				rpolicyID = prompt.GetReplicationPolicyFromUser()
			}
			response, err := api.StartReplication(rpolicyID)
			if err != nil {
				return fmt.Errorf("failed to start replication: %v", utils.ParseHarborErrorMsg(err))
			}
			fmt.Printf("Repliation started successfully with ID: %s\n", response.Location)
			return nil
		},
	}

	return cmd
}
