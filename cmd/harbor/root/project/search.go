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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SearchProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "search project based on their names",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("Starting project search command")
			projects, err := api.SearchProject(args[0])
			if err != nil {
				return fmt.Errorf("failed to get projects: %v", utils.ParseHarborErrorMsg(err))
			}
			log.Debugf("Found %d projects", len(projects.Payload.Project))
			if len(projects.Payload.Project) == 0 {
				return fmt.Errorf("No projects found with name similar to : %s", args[0])
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(projects, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				list.SearchProjects(projects.Payload.Project)
			}
			return nil
		},
	}
	return cmd
}
