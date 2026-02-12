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
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/repository/search"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SearchRepoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "search repository based on their names",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := api.SearchRepository(args[0])
			if err != nil {
				return fmt.Errorf("failed to get repositories: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(repo, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				search.SearchRepositories(repo.Payload.Repository)
			}
			return nil
		},
	}
	return cmd
}
