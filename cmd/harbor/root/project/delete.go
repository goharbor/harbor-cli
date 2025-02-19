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
package project

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DeleteProjectCommand creates a new `harbor delete project` command
func DeleteProjectCommand() *cobra.Command {
	var forceDelete bool
	var useProjectID bool

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete project by name or ID",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				err = api.DeleteProject(args[0], forceDelete, useProjectID)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				err = api.DeleteProject(projectName, forceDelete, false)
			}
			if err != nil {
				log.Errorf("failed to delete project: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&forceDelete, "force", false, "Deletes all repositories and artifacts within the project")
	flags.BoolVar(&useProjectID, "projectID", false, "If set, treats the provided argument as a project ID instead of a project name")

	return cmd
}
