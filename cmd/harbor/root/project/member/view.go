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
	"github.com/goharbor/harbor-cli/pkg/views/member/view"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ViewMemberCommand creates a new `harbor project member view` command
func ViewMemberCommand() *cobra.Command {
	var opts api.GetMemberOptions
	var isID bool

	cmd := &cobra.Command{
		Use:     "view [projectName] [memberID]",
		Short:   "get project member information",
		Example: "  harbor project member view my-project 5\n  harbor project member view my-project 5 --wide",
		Args:    cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) == 1 {
				opts.ProjectNameOrID = args[0]
			} else if len(args) == 2 {
				opts.ProjectNameOrID = args[0]
				opts.ID, _ = strconv.ParseInt(args[1], 0, 64)
			}

			if opts.ProjectNameOrID == "" {
				opts.ProjectNameOrID, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
			}

			if opts.ID == 0 {
				opts.ID = prompt.GetMemberIDFromUser(opts.ProjectNameOrID, "")
				if opts.ID == 0 {
					return fmt.Errorf("no members found in project")
				}
			}

			opts.XIsResourceName = !isID

			member, err := api.GetMember(opts)
			if err != nil {
				return fmt.Errorf("failed to get project member info: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(member.Payload, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				wideFlag := viper.GetBool("wide")
				view.ViewMember(member.Payload, wideFlag)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "parses projectName as an ID")
	return cmd
}
