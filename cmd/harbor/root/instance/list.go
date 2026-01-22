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
package instance

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/instance/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListInstanceCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all preheat provider instances in Harbor",
		Long: `List all preheat provider instances registered in Harbor. You can paginate the results, 
filter them using a query string, and sort them in ascending or descending order. 
This command provides an easy way to view all instances along with their details.`,
		Example: `  harbor-cli instance list --page 1 --page-size 10
  harbor-cli instance list --query "name=my-instance" --sort "asc"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}
			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

			instance, err := api.ListInstance(opts)

			if err != nil {
				log.Fatalf("failed to get instance list: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(instance, FormatFlag)
				if err != nil {
					log.Errorf("Failed to print config: %v", err)
				}
			} else {
				list.ListInstance(instance.Payload)
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
