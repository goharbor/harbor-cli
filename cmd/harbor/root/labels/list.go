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
	var (
		opts        api.ListFlags
		projectName string
		isGlobal    bool
		// For querying, opts.Q
		fuzzy  []string
		match  []string
		ranges []string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list labels",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}

			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

			// Defining ProjectID & Scope based on user inputs
			if isGlobal {
				opts.Scope = "g"
			} else if projectName != "" {
				id, err := api.GetProjectIDFromName(projectName)
				if err != nil {
					return err
				}

				opts.ProjectID = id
				opts.Scope = "p"
			} else if opts.ProjectID != 0 {
				opts.Scope = "p"
			} else {
				opts.Scope = "g"
			}

			if len(fuzzy) != 0 || len(match) != 0 || len(ranges) != 0 { // Only Building Query if a param exists
				q, qErr := utils.BuildQueryParam(fuzzy, match, ranges,
					[]string{"name", "id", "label_id", "creation_time", "owner_id", "color", "description"},
				)
				if qErr != nil {
					return qErr
				}

				opts.Q = q
			}

			label, err := api.ListLabel(opts)
			if err != nil {
				log.Fatalf("failed to get label list: %v", err)
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
	flags.StringVarP(&projectName, "project", "p", "", "project name when query project labels")
	flags.Int64VarP(&opts.ProjectID, "project-id", "i", 0, "project ID when query project labels")
	flags.BoolVarP(&isGlobal, "global", "", false, "whether to list global or project scope labels. (default scope is global)")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the label list in ascending or descending order")
	flags.StringSliceVar(&fuzzy, "fuzzy", nil, "Fuzzy match filter (key=value)")
	flags.StringSliceVar(&match, "match", nil, "exact match filter (key=value)")
	flags.StringSliceVar(&ranges, "range", nil, "range filter (key=min~max)")

	return cmd
}
