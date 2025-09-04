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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/member/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListMemberCommand creates a new `harbor member list` command
func ListMemberCommand() *cobra.Command {
	var opts api.ListMemberOptions

	cmd := &cobra.Command{
		Use:     "list [projectName]",
		Short:   "list members in a project",
		Long:    "list members in a project by projectName",
		Example: "  harbor project member list my-project",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) > 0 {
				opts.ProjectNameOrID = args[0]
			} else {
				opts.ProjectNameOrID, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
			}

			members, err := api.ListMember(opts)
			if err != nil {
				return fmt.Errorf("failed to get members list: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag == "json" {
				utils.PrintPayloadInJSONFormat(members)
				return nil
			} else if FormatFlag == "yaml" {
				utils.PrintPayloadInYAMLFormat(members)
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
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.EntityName, "name", "n", "", "Member Name to search")

	return cmd
}
