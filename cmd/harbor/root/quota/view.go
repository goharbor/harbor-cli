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
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/quota/list"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var quota *models.Quota

			// get quota id with quota
			quota, err = GetQuotaFromUser(args, opts)
			if err != nil {
				return fmt.Errorf("failed to get quota: %w", err)
			}
			quotas := []*models.Quota{quota}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(quota, FormatFlag)
				if err != nil {
					return fmt.Errorf("failed to get quota list: %w", err)
				}
			} else {
				list.ListQuotas(quotas)
			}
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Reference, "project-name", "", "", "Get quota by project-name")
	flags.StringVarP(&opts.ReferenceID, "project-id", "", "", "Get quota by project ID")

	return cmd
}
