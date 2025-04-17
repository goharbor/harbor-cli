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
	"fmt"
	"sync"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			var wg sync.WaitGroup
			var mu sync.Mutex

			successfulDeletes := []string{}
			failedDeletes := map[string]string{}

			if len(args) > 0 {
				log.Debugf("Deleting %d projects from args...", len(args))
				for _, projectName := range args {
					pn := projectName
					log.Debugf("Initiating delete for project: %s", pn)
					wg.Add(1)
					go func(projectName string) {
						defer wg.Done()
						log.Debugf("Deleting project '%s' with force=%v", projectName, forceDelete)
						if err := api.DeleteProject(projectName, forceDelete, false); err != nil {
							mu.Lock()
							failedDeletes[projectName] = utils.ParseHarborError(err)
							mu.Unlock()
						} else {
							mu.Lock()
							successfulDeletes = append(successfulDeletes, projectName)
							mu.Unlock()
						}
					}(pn)
				}
			} else {
				log.Debug("No arguments provided. Prompting user for project name.")
				projectName := prompt.GetProjectNameFromUser()
				log.Debugf("User input project: %s", projectName)
				log.Debugf("Deleting project '%s' with force=%v", projectName, forceDelete)
				if err := api.DeleteProject(projectName, forceDelete, false); err != nil {
					return fmt.Errorf("failed to delete project: %v", utils.ParseHarborError(err))
				}
				fmt.Printf("Project '%s' deleted successfully\n", projectName)
				return nil
			}

			wg.Wait()

			if len(successfulDeletes) > 0 {
				fmt.Println("Successfully deleted projects:")
				for _, name := range successfulDeletes {
					fmt.Printf("  - %s\n", name)
				}
			}

			if len(failedDeletes) > 0 {
				fmt.Println("Failed to delete projects:")
				for name, reason := range failedDeletes {
					fmt.Printf("  - %s: %s\n", name, reason)
				}
				return fmt.Errorf("failed to delete %d project(s)", len(failedDeletes))
			}

			log.Debug("All requested projects deleted successfully.")
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&forceDelete, "force", false, "Deletes all repositories and artifacts within the project")

	return cmd
}
