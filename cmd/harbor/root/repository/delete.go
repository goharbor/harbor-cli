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
package repository

import (
	"sync"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RepoDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a repository",
		Example: `  harbor repository delete [project_name]/[repository_name]`,
		Long:    `Delete a repository within a project in Harbor`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				var wg sync.WaitGroup
				errChan := make(chan error, len(args))
				for _, arg := range args {
					wg.Add(1)
					go func(arg string) {
						defer wg.Done()
						projectName, repoName := utils.ParseProjectRepo(arg)
						err := api.RepoDelete(projectName, repoName)
						errChan <- err
						if err != nil {
							log.Errorf("failed to delete repository %v: %v", arg, err)
						}
					}(arg)
				}
				wg.Wait()
				close(errChan)
				var countDeleted int
				var countErr int
				var finalErr error
				for err := range errChan {
					if err != nil {
						countErr++
						finalErr = err
					} else {
						countDeleted++
					}
				}
				err = finalErr
				log.Infof("deleted %d repositories, %d failed", countDeleted, countErr)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				err = api.RepoDelete(projectName, repoName)
			}
			if err != nil {
				log.Errorf("failed to delete repository: %v", err)
			}
		},
	}
	return cmd
}
