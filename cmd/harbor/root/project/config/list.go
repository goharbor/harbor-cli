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
package config

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ListConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list [NAME|ID]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name")
			} else {
				projectNameOrID := args[0]

				response, err := api.ListConfig(isID, projectNameOrID)
				if err != nil {
					log.Errorf("failed to view metadata: %v", err)
				} else {
					utils.PrintPayloadInJSONFormat(response.Payload)
				}
			}
		},
	}
	return cmd
}
