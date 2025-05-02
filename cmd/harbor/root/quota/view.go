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
package quota

import (
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/quota/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// View a specified quota
func ViewQuotaCommand() *cobra.Command {
	var opts api.ListQuotaFlags
	cmd := &cobra.Command{
		Use:   "view [quotaID]",
		Short: "get quota by quota ID",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var quota *models.Quota

			if len(args) > 0 {
				quotaID, _ := strconv.ParseInt(args[0], 10, 64)
				quota, err = api.GetQuota(int64(quotaID))
				if err != nil {
					log.Errorf("failed to get Quota: %v", err)
				}
			} else if opts.Reference != "" {
				project, err := api.GetProject(opts.Reference, false)
				if err != nil {
					log.Errorf("failed to get project: %v", err)
				}
				projectID := project.Payload.ProjectID
				quota, err = api.GetQuotaByRef(int64(projectID))
				if err != nil {
					log.Errorf("failed to get quota: %v", err)
					return
				}
			} else if opts.ReferenceID != "" {
				projectID, err := strconv.ParseInt(opts.ReferenceID, 10, 64)
				if err != nil {
					log.Errorf("invalid projectID: %v", err)
					return
				}
				quota, err = api.GetQuotaByRef(projectID)
				if err != nil {
					log.Errorf("failed to get quota: %v", err)
					return
				}
			} else {
				quotaID := prompt.GetQuotaIDFromUser()
				if quotaID == 0 {
					log.Errorf("failed to get quotaID from user")
					return
				}
				quota, err = api.GetQuota(quotaID)
			}

			if err != nil {
				log.Errorf("failed to get project: %v", err)
			}

			quotas := []*models.Quota{quota}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(quota, FormatFlag)
				if err != nil {
					log.Errorf("failed to get quota list: %v", err)
					return
				}
			} else {
				list.ListQuotas(quotas)
			}

		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Reference, "project-name", "", "", "Get quota by project-name")
	flags.StringVarP(&opts.ReferenceID, "project-id", "", "", "Get quota by project ID")

	return cmd
}
