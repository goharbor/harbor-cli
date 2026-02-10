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
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func RepoDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a repository",
		Example: `  harbor repository delete [project_name]/[repository_name]`,
		Long:    `Delete a repository within a project in Harbor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName string
			var repoName string
			if len(args) > 0 {
				projectName, repoName, err = utils.ParseProjectRepo(args[0])
				if err != nil {
					return fmt.Errorf("failed to parse project/repo: %v", err)
				}
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
				repoName = prompt.GetRepoNameFromUser(projectName)
			}
			err = api.RepoDelete(projectName, repoName, false)
			if err != nil {
				return fmt.Errorf("failed to delete repository: %v", err)
			}
			return nil
		},
	}
	return cmd
}
