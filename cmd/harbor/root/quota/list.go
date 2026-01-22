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
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/quota/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Lists the Quotas specified for each project
func ListQuotaCommand() *cobra.Command {
	var opts api.ListQuotaFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list quotas",
		Long:  "list quotas specified for each project",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.PageSize < 0 {
				log.Errorf("page size must be greater than or equal to 0")
				return
			}

			if opts.PageSize > 100 {
				log.Errorf("page size should be less than or equal to 100")
				return
			}

			quota, err := api.ListQuota(opts)
			if err != nil {
				log.Errorf("failed to get quota list: %v", err)
				return
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(quota, FormatFlag)
				if err != nil {
					log.Errorf("failed to get quota list: %v", err)
					return
				}
			} else {
				list.ListQuotas(quota.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 0, "Size of per page (use 0 to fetch all)")
	flags.StringVarP(&opts.Reference, "ref", "", "", "Reference type of quota")
	flags.StringVarP(&opts.ReferenceID, "refid", "", "", "Reference ID of quota")
	flags.StringVarP(
		&opts.Sort,
		"sort",
		"",
		"",
		"Sort the resource list in ascending or descending order",
	)

	return cmd
}
