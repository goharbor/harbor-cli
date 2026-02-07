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

	proj "github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	auditLog "github.com/goharbor/harbor-cli/pkg/views/project/logs"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LogsProjectCommmand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "get project logs",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("Starting execution of 'logs' command")
			var err error
			var resp *proj.GetLogExtsOK
			var projectName string

			if len(args) > 0 {
				projectName = args[0]
				log.Debugf("Project name provided as argument: %s", projectName)
			} else {
				log.Debug("No project name argument provided, prompting user...")
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
				log.Debugf("Project name received from prompt: %s", projectName)
			}

			log.Debugf("Checking if project '%s' exists...", projectName)
			_, err = api.GetProject(projectName, false)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "404" {
					return fmt.Errorf("project %s does not exist", projectName)
				}
				return fmt.Errorf("failed to verify project: %v", utils.ParseHarborErrorMsg(err))
			}

			log.Debugf("Fetching logs for project: %s", projectName)
			resp, err = api.LogsProject(projectName)
			if err != nil {
				return fmt.Errorf("failed to get project logs: %v", utils.ParseHarborErrorMsg(err))
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				log.WithField("output_format", formatFlag).Debug("Output format selected")
				err = utils.PrintFormat(resp.Payload, formatFlag)
				if err != nil {
					return err
				}
			} else {
				log.Debug("Listing project logs using default view")
				auditLog.LogsProject(resp.Payload)
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
