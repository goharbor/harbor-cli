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
package labels

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/label/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListLabelCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list labels",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

			if !cmd.Flag(("page-size")).Changed {
				if defaultPageSize, ok := utils.GetDefaultPageSize(); ok {
					opts.PageSize = defaultPageSize
				}
			}

			label, err := api.ListLabel(opts)
			if err != nil {
				log.Fatalf("failed to get label list: %v", err)
			}
			if len(label.Payload) == 0 {
				log.Info("No labels found")
				return nil
			}
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(label, formatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.ListLabels(label.Payload)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 20, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "default(global).'p' for project labels.Query scope of the label")
	flags.Int64VarP(&opts.ProjectID, "projectid", "i", 1, "project ID when query project labels")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the label list in ascending or descending order")

	return cmd
}
