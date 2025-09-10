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
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	view "github.com/goharbor/harbor-cli/pkg/views/member/view"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewGetRegistryCommand creates a new `harbor get registry` command
func ViewMemberCommand() *cobra.Command {
	var opts api.GetMemberOptions
	var isID bool

	cmd := &cobra.Command{
		Use:     "view [ProjectName] [member ID]",
		Short:   "get project member details",
		Long:    "get member details by MemberID",
		Example: "  harbor project member view my-project [memberID]",
		Args:    cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) == 1 {
				opts.ProjectNameOrID = args[0]
			} else if len(args) == 2 {
				opts.ID, _ = strconv.ParseInt(args[1], 0, 64)
				opts.ProjectNameOrID = args[0]
			} else if opts.ProjectNameOrID == "" || opts.ID == 0 {
				if opts.ProjectNameOrID == "" {
					opts.ProjectNameOrID, err = prompt.GetProjectNameFromUser()
					if err != nil {
						return fmt.Errorf("failed to get project name: %v", err)
					}
				}
				if opts.ID == 0 {
					opts.ID = prompt.GetMemberIDFromUser(opts.ProjectNameOrID)
				}
			}

			// when set true parses projectNameOrID as projectName
			// else it parses as an integer ID
			opts.XIsResourceName = !isID

			member, err := api.GetMember(opts)
			if err != nil {
				return fmt.Errorf("failed to get members list: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			VerboseFlag := viper.GetBool("verbose")
			if FormatFlag == "json" {
				err = utils.PrintFormat(member, FormatFlag)
				return err
			} else if VerboseFlag {
				view.ViewMember(member.Payload, true)
			} else {
				view.ViewMember(member.Payload, false)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "Parses projectName as an ID")
	flags.Int64VarP(&opts.ID, "member-id", "", 0, "Member ID")
	flags.StringVarP(&opts.ProjectNameOrID, "projectname", "p", "", "Project Name")
	return cmd
}
