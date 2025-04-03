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
	"sync"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DeleteProjectCommand creates a new `harbor project delete` command
func DeleteProjectCommand() *cobra.Command {
	var forceDelete bool
	var projectID string

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete project by name or ID",
		Example: "harbor project delete [projectname] or harbor project delete --project-id [projectid]",
		Long:    "Delete project by name or ID. If no arguments are provided, it will prompt for the project name. Use --project-id to specify the project ID directly. The --force flag will delete all repositories and artifacts within the project.",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			errChan := make(chan error, len(args))

			if projectID != "" {
				wg.Add(1)
				go func(id string) {
					defer wg.Done()
					if err := api.DeleteProject(id, forceDelete, true); err != nil {
						errChan <- err
					}
				}(projectID)
			} else if len(args) > 0 {
				for _, projectName := range args {
					wg.Add(1)
					go func(name string) {
						defer wg.Done()
						if err := api.DeleteProject(name, forceDelete, false); err != nil {
							errChan <- err
						}
					}(projectName)
				}
			} else {
				projectName := prompt.GetProjectNameFromUser()
				if err := api.DeleteProject(projectName, forceDelete, false); err != nil {
					log.Errorf("failed to delete project: %v", err)
				}
			}

			go func() {
				wg.Wait()
				close(errChan)
			}()

			var finalErr error
			for err := range errChan {
				if finalErr == nil {
					finalErr = err
				} else {
					log.Errorf("Error: %v", err)
				}
			}
			if finalErr != nil {
				log.Errorf("failed to delete some projects: %v", finalErr)
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&forceDelete, "force", false, "Deletes all repositories and artifacts within the project")
	flags.StringVar(&projectID, "project-id", "", "Specify project ID instead of project name")

	return cmd
}
