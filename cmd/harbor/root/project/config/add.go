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
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func AddConfigCommand() *cobra.Command {
	var err error
	var projectNameorID string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add [NAME|ID] ...[KEY]:[VALUE]",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				projectNameorID, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
				isID = false
			} else {
				projectNameorID = args[0]
			}
			metadata := make(map[string]string)
			for i := 1; i < len(args); i++ {
				keyValue := args[i]
				keyValueArray := strings.Split(keyValue, ":")
				if len(keyValueArray) == 2 {
					metadata[keyValueArray[0]] = keyValueArray[1]
				} else {
					return fmt.Errorf("Please provide metadata in the format key:value")

				}
			}
			err = api.AddConfig(isID, projectNameorID, metadata)
			if err != nil {
				return fmt.Errorf("failed to add metadata: %v", err)
			}
			return nil
		},
	}
	return cmd

}
