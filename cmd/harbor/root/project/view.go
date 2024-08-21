// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package project

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// GetProjectCommand creates a new `harbor get project` command
func ViewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "view [NAME|ID]",
		Short: "get project by name or id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				err = api.GetProject(args[0])
			} else {
				projectName := prompt.GetProjectNameFromUser()
				err = api.GetProject(projectName)
			}

			if err != nil {
				log.Errorf("failed to get project: %v", err)
			}

		},
	}

	return cmd
}
