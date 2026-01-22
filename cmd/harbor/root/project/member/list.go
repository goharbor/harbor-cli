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

package member

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/member/list"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListMemberCommand creates a new `harbor member list` command
func ListMemberCommand() *cobra.Command {
	var opts api.ListMemberOptions
	var isID bool
	var searchQuery string

	cmd := &cobra.Command{
		Use:     "list [projectName]",
		Short:   "list members in a project",
		Long:    "list members in a project by projectName",
		Example: "  harbor project member list my-project",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}
			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}

			var err error
			if len(args) > 0 {
				opts.ProjectNameOrID = args[0]
			} else {
				opts.ProjectNameOrID, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
			}

			// when set true parses projectNameOrID as projectName
			// else it parses as an integer ID
			opts.XIsResourceName = !isID

			members, err := api.ListMember(opts)
			if err != nil {
				return fmt.Errorf("failed to get members list: %v", err)
			}

			if searchQuery != "" && opts.EntityName == "" {
				set := make([]string, 0, len(members.Payload))
				for _, v := range members.Payload {
					set = append(set, v.EntityName)
				}

				matches := fuzzy.Find(searchQuery, set)

				results := make([]*models.ProjectMemberEntity, 0)
				for _, v := range matches {
					results = append(results, members.Payload[v.Index])
				}

				members.Payload = results
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(members, FormatFlag)
				if err != nil {
					return err
				}

				return nil
			}

			VerboseFlag := viper.GetBool("verbose")

			if VerboseFlag {
				list.ListMembers(members.Payload, true)
			} else {
				list.ListMembers(members.Payload, false)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "Parses projectName as an ID")
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&searchQuery, "fuzzy", "f", "", "Fuzzy search for member with name")
	flags.StringVarP(&opts.EntityName, "search", "s", "", "Search for member with name")

	return cmd
}
