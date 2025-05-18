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
	"github.com/goharbor/harbor-cli/pkg/views/label/update"
	"github.com/spf13/cobra"
)

func UpdateLableCommand() *cobra.Command {
	opts := &models.Label{}

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "update label",
		Example: "harbor label update [labelname]",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var labelId int64
			updateflags := api.ListFlags{
				Scope:     opts.Scope,
				ProjectID: opts.ProjectID,
			}

			if len(args) > 0 {
				labelId, err = api.GetLabelIdByName(args[0], updateflags)
			} else {
				labelId, err = prompt.GetLabelIdFromUser(updateflags)
			}
			if err != nil {
				return fmt.Errorf("failed to parse label id: %v", err)
			}

			existingLabel := api.GetLabel(labelId)
			if existingLabel == nil {
				return fmt.Errorf("label is not found")
			}
			updateView := &models.Label{
				Name:        existingLabel.Name,
				Color:       existingLabel.Color,
				Description: existingLabel.Description,
				Scope:       existingLabel.Scope,
			}

			flags := cmd.Flags()
			if flags.Changed("name") {
				updateView.Name = opts.Name
			}
			if flags.Changed("color") {
				updateView.Color = opts.Color
			}
			if flags.Changed("description") {
				updateView.Description = opts.Description
			}
			if flags.Changed("scope") {
				updateView.Scope = opts.Scope
			}

			update.UpdateLabelView(updateView)
			err = api.UpdateLabel(updateView, labelId)
			if err != nil {
				return fmt.Errorf("failed to update label: %v", err)
			}
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "Name of the label")
	flags.StringVarP(&opts.Color, "color", "", "", "Color of the label.color is in hex value")
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "Scope of the label. eg- g(global), p(specific project)")
	flags.StringVarP(&opts.Description, "description", "d", "", "Description of the label")
	flags.Int64VarP(&opts.ProjectID, "project", "i", 0, "Description of the label")

	return cmd
}
