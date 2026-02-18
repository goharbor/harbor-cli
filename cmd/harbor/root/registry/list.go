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
package registry

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/registry/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewListRegistryCommand creates a new `harbor list registry` command
func ListRegistryCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list registry",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}

			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}
			registry, err := api.ListRegistries(opts)

			if err != nil {
				return fmt.Errorf("failed to get projects list: %v", err)
			}
			if len(registry.Payload) == 0 {
				log.Info("No registries found")
				return nil
			}
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(registry, formatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.ListRegistry(registry.Payload)
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
