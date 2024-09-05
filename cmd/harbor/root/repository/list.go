// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repository

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/repository/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ListRepositoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list repositories within a project",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var resp repository.ListRepositoriesOK

			if len(args) > 0 {
				resp, err = api.ListRepository(args[0])
			} else {
				projectName := prompt.GetProjectNameFromUser()
				resp, err = api.ListRepository(projectName)
			}

			if err != nil {
				log.Errorf("failed to list repositories: %v", err)
			}

			list.ListRepositories(resp.Payload)

		},
	}

	return cmd
}
