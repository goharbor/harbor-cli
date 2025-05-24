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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func DeleteLabelCommand() *cobra.Command {
	var opts models.Label
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete label",
		Example: "harbor label delete [labelname]",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var labelId int64
			deleteView := &api.ListFlags{
				Scope: opts.Scope,
			}

			if len(args) > 0 {
				labelId, _ = api.GetLabelIdByName(args[0])
			} else {
				labelList, err := api.ListLabel(*deleteView)
				if err != nil {
					return fmt.Errorf("failed to get label list: %v", utils.ParseHarborErrorMsg(err))
				}
				labelId = prompt.GetLabelIdFromUser(labelList.Payload)
			}
			err = api.DeleteLabel(labelId)
			if err != nil {
				return fmt.Errorf("failed to delete label: %v", utils.ParseHarborErrorMsg(err))
			}
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "default(global).'p' for project labels.Query scope of the label")

	return cmd
}
