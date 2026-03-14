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
	"github.com/goharbor/harbor-cli/pkg/views/repository/update"
	"github.com/spf13/cobra"
)

func UpdateRepositoryCommand() *cobra.Command {
	var description string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a repository",
		Long: `Update the description of a repository.

This command updates the description associated with a repository
within a Harbor project.

Examples:
  # Update repository description using project/repository format
  	harbor repository update library/nginx --description "Official nginx image"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName string
			var repoName string

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

			existingRepo, err := api.RepoView(projectName, repoName)
			if err != nil {
				return fmt.Errorf("failed to get existing repository: %w", err)
			}

			if cmd.Flags().Changed("description") {
				err = api.UpdateRepository(projectName, repoName, description, existingRepo.Payload)
				if err != nil {
					return fmt.Errorf("failed to update repository: %w", err)
				}
			} else {
				updatedDescription, err := update.UpdateRepositoryView(existingRepo.Payload)
				if err != nil {
					return fmt.Errorf("update cancelled or failed: %w", err)
				}

				err = api.UpdateRepository(projectName, repoName, updatedDescription, existingRepo.Payload)
				if err != nil {
					return fmt.Errorf("failed to update repository: %w", err)
				}
			}

			fmt.Printf("Repository %s/%s updated successfully\n", projectName, repoName)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&description, "description", "d", "", "Repository description")

	return cmd
}
