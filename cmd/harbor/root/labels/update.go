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
	"github.com/goharbor/harbor-cli/pkg/views/label/update"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func UpdateLableCommand() *cobra.Command {
	opts := &models.Label{}
	var projectName string
	var isGlobal bool

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "update label",
		Example: "harbor label update [labelname]",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var labelId int64

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

			existingLabel, err := api.GetLabel(labelId)
			if err != nil {
				return fmt.Errorf("failed to get label: %v", utils.ParseHarborErrorMsg(err))
			}
			updateView := &models.Label{
				Name:        existingLabel.Name,
				Color:       existingLabel.Color,
				Description: existingLabel.Description,
				Scope:       existingLabel.Scope,
			}

			// Only open the interactive form when no update flags were given,
			// so the command can be used non-interactively (scripts, CI).
			if !applyLabelUpdateFlags(cmd.Flags(), opts, updateView) {
				update.UpdateLabelView(updateView)
			} else if err := validateLabelUpdate(updateView); err != nil {
				return err
			}
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
	flags.StringVarP(&projectName, "project", "p", "", "project name when query project labels")
	flags.BoolVarP(&isGlobal, "global", "", false, "whether to list global or project scope labels. (default scope is global)")
	flags.Int64VarP(&opts.ProjectID, "project-id", "i", 0, "project ID when query project labels")
	flags.StringVarP(&opts.Description, "description", "d", "", "Description of the label")

	return cmd
}

// applyLabelUpdateFlags copies the explicitly set update flags onto updateView
// and reports whether any of them were set.
func applyLabelUpdateFlags(flags *pflag.FlagSet, opts *models.Label, updateView *models.Label) bool {
	changed := false
	if flags.Changed("name") {
		updateView.Name = opts.Name
		changed = true
	}
	if flags.Changed("color") {
		updateView.Color = opts.Color
		changed = true
	}
	if flags.Changed("description") {
		updateView.Description = opts.Description
		changed = true
	}
	return changed
}

// validateLabelUpdate enforces the same rules as the interactive form for
// values supplied through flags.
func validateLabelUpdate(updateView *models.Label) error {
	if updateView.Name == "" {
		return fmt.Errorf("label name cannot be empty")
	}
	if updateView.Color == "" {
		return fmt.Errorf("label color cannot be empty")
	}
	return nil
}
