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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/summary"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getProjectSummaryFunc = api.GetProjectSummary

func SummaryCommand() *cobra.Command {
	var isID bool
	cmd := &cobra.Command{
		Use:   "summary [PROJECT_NAME|PROJECT_ID]",
		Short: "Get summary of a project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			var err error

			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			log.Debugf("Fetching project metadata: %s", projectName)
			projectData, err := api.GetProject(projectName, isID)
			if err != nil {
				return handleProjectError(err, projectName, "get details of")
			}

			log.Debugf("Fetching project summary: %s", projectName)
			projectSummary, err := getProjectSummaryFunc(projectName, isID)
			if err != nil {
				return handleProjectError(err, projectName, "get summary of")
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				log.WithField("output_format", FormatFlag).Debug("Output format selected")
				if FormatFlag == "csv" {
					type combinedCSV struct {
						Project *models.Project        `json:"project" csv:"project"`
						Summary *models.ProjectSummary `json:"summary" csv:"summary"`
					}
					combined := []combinedCSV{{
						Project: projectData.Payload,
						Summary: projectSummary.Payload,
					}}
					err = utils.PrintFormat(combined, FormatFlag)
				} else {
					combined := map[string]interface{}{
						"project": projectData.Payload,
						"summary": projectSummary.Payload,
					}
					err = utils.PrintFormat(combined, FormatFlag)
				}
				if err != nil {
					return err
				}
			} else {
				log.Debug("Showing project summary using default view")
				err = summary.ViewProjectSummary(projectData.Payload, projectSummary.Payload)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&isID, "id", "i", false, "Identify project by ID instead of name")

	return cmd
}

func handleProjectError(err error, name string, action string) error {
	errorCode := utils.ParseHarborErrorCode(err)
	switch errorCode {
	case "401":
		return fmt.Errorf("unauthorized: please login to Harbor")
	case "403":
		return fmt.Errorf("forbidden: you do not have permission to %s project %s", action, name)
	case "404":
		return fmt.Errorf("project %s does not exist", name)
	case "500":
		return fmt.Errorf("internal server error: please contact your administrator")
	default:
		return fmt.Errorf("failed to %s project %s: %v", action, name, utils.ParseHarborErrorMsg(err))
	}
}
