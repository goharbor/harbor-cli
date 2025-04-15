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

// DeleteProjectCommand creates a new `harbor delete project` command
func DeleteProjectCommand() *cobra.Command {
	var forceDelete bool
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete project by name or id",
		Example: `  harbor project delete [projectname]`,
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			errChan := make(
				chan error,
				len(args),
			) // Channel to collect errors

			if len(args) > 0 {
				for _, arg := range args {
					wg.Add(1)
					go func(projectName string) {
						defer wg.Done()
						if err := api.DeleteProject(projectName, forceDelete); err != nil {
							errChan <- err
						}
					}(arg)
				}
			} else {
				projectName := prompt.GetProjectNameFromUser()
				err := api.DeleteProject(projectName, forceDelete)
				if err != nil {
					log.Errorf("failed to delete project: %v", err)
				}
			}

			// Wait for all goroutines to finish
			go func() {
				wg.Wait()
				close(errChan)
			}()

			// Collect and handle errors
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

	return cmd
}
