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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/repository/view"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RepoViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "Get repository information",
		Example: `  harbor repo view <project_name>/<repo_name>`,
		Long:    `Get information of a particular repository in a project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName, repoName string
			var repo *repository.GetRepositoryOK

			if len(args) > 0 {
				projectName, repoName, err = utils.ParseProjectRepo(args[0])
				if err != nil {
					return fmt.Errorf("failed to parse project/repo: %w", err)
				}
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %w", err)
				}
				repoName = prompt.GetRepoNameFromUser(projectName)
			}

			repo, err = api.RepoView(projectName, repoName)
			if err != nil {
				return fmt.Errorf("failed to get repository information: %w", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(repo, FormatFlag)
				if err != nil {
					return fmt.Errorf("failed to print format: %w", err)
				}
			} else {
				view.ViewRepository(repo.Payload)
			}
			return nil
		},
	}

	return cmd
}
