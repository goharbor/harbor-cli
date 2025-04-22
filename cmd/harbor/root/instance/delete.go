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
package instance

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteInstanceCommand() *cobra.Command {
	var instanceID int64
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a preheat provider instance by its name or ID",
		Long: `Delete a preheat provider instance from Harbor. You can specify the instance name or ID directly as an argument.
If no argument is provided, you will be prompted to select an instance from a list of available instances.`,
		Example: `  harbor-cli instance delete my-instance
  harbor-cli instance delete 12345`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var instanceName string
			
			if instanceID != -1 {
				instanceName, err = api.GetInstanceNameByID(instanceID)
				if err != nil {
					log.Errorf("%v", err)
					return
				}
			} else if len(args) > 0 {
				instanceName = args[0]
			} else {
				instanceName = prompt.GetInstanceFromUser()
			}
			err = api.DeleteInstance(instanceName)
			if err != nil {
				log.Errorf("failed to delete instance: %v", err)
			}
		},
	}
	cmd.Flags().Int64VarP(&instanceID, "id", "i", -1, "ID of the instance to delete")
	return cmd
}
