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
	"strings"

	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update [NAME|ID] [KEY] ...[KEY]:[VALUE]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name, the meta name and metadata")
			} else if len(args) == 1 {
				fmt.Println("Please provide the meta name and metadata")
			} else if len(args) == 2 {
				fmt.Println("Please provide metadata")
			} else {
				projectNameOrID := args[0]
				metaName := args[1]
				metadata := make(map[string]string)
				for i := 2; i < len(args); i++ {
					keyValue := args[i]
					keyValueArray := strings.Split(keyValue, ":")
					if len(keyValueArray) == 2 {
						metadata[keyValueArray[0]] = keyValueArray[1]
					} else {
						fmt.Println("Please provide metadata in the format key:value")
						return
					}
				}
				err := api.UpdateConfig(isID, projectNameOrID, metaName, metadata)
				if err != nil {
					log.Errorf("failed to view metadata: %v", err)
				}
			}
		},
	}
	return cmd
}
