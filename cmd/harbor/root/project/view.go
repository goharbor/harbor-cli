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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ViewCommand() *cobra.Command {
	var isID bool
	cmd := &cobra.Command{
		Use:   "view [NAME|ID]",
		Short: "get project by name or id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var projectNameOrID string
			var project *project.GetProjectOK

			if len(args) > 0 {
				projectNameOrID = args[0]
			} else {
				projectNameOrID = prompt.GetProjectNameFromUser()
			}

			project, err = api.GetProject(projectNameOrID, false)
			if err != nil {
				log.Errorf("failed to get project: %v", err)
				return
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(project, FormatFlag)
				if err != nil {
					log.Error(err)
					return
				}
			} else {
				view.ViewProjects(project.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&isID, "id", false, "Get project by id")

	return cmd
}
