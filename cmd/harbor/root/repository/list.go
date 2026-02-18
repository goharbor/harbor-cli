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
	"github.com/goharbor/harbor-cli/pkg/views/repository/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListRepositoryCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list repositories within a project",
		Example: `  harbor repo list <project_name>`,
		Long:    `Get information of all repositories in a project`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}

			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}
			var err error
			var repos repository.ListRepositoriesOK
			var projectName string

			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			repos, err = api.ListRepository(projectName, false, opts)
			if err != nil {
				return fmt.Errorf("failed to list repositories: %v", err)
			}
			if len(repos.Payload) == 0 {
				log.Info("No repositories found")
				return nil
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(repos, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.ListRepositories(repos.Payload)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")

	return cmd
}
