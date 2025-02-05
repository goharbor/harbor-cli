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
		Use:   "delete",
		Short: "delete project by name or id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				var wg sync.WaitGroup
				errChan := make(chan error, len(args))
				for _, arg := range args {
					wg.Add(1)
					go func(arg string) {
						defer wg.Done()
						err := api.DeleteProject(arg, forceDelete)
						errChan <- err
						if err != nil {
							log.Errorf("failed to delete project %v: %v", arg, err)
						}
					}(arg)
				}
				wg.Wait()
				close(errChan)
				var countdeleted int
				var counterr int
				var finalerr error
				for err := range errChan {
					if err != nil {
						counterr++
						finalerr = err
					} else {
						countdeleted++
					}
				}
				err = finalerr
				log.Infof("deleted %d projects, %d failed", countdeleted, counterr)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				err = api.DeleteProject(projectName, forceDelete)
			}
			if err != nil {
				log.Errorf("failed to delete project: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&forceDelete, "force", false, "Deletes all repositories and artifacts within the project")

	return cmd
}
