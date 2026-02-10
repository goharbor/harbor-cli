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
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/label/create"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateLabelCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create label",
		Long:    "create label in harbor",
		Example: "harbor label create",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			createView := &create.CreateView{
				Name:        opts.Name,
				Color:       opts.Color,
				Scope:       opts.Scope,
				Description: opts.Description,
				ProjectID:   opts.ProjectID,
			}
			if opts.Name != "" && opts.Scope != "" {
				if opts.Scope == "p" && opts.ProjectID == 0 {
					return fmt.Errorf("project ID is required when scope is 'p' (project-specific). Use --project flag to specify the project ID")
				}
				if opts.Scope == "p" {
					_, err := api.GetProject(strconv.FormatInt(opts.ProjectID, 10), true)
					if err != nil {
						return fmt.Errorf("project with ID %d does not exist", opts.ProjectID)
					}
				}
				err = api.CreateLabel(opts)
			} else {
				flags := cmd.Flags()
				err = createLabelView(createView, flags)
			}

			if err != nil {
				return fmt.Errorf("failed to create label: %v", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "Name of the label")
	flags.StringVarP(&opts.Color, "color", "", "#FFFFFF", "Color of the label.color is in hex value")
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "Scope of the label. eg- g(global), p(specific project)")
	flags.Int64VarP(&opts.ProjectID, "project", "i", 0, "Id of the project when scope is p")
	flags.StringVarP(&opts.Description, "description", "d", "", "Description of the label")

	return cmd
}

func createLabelView(createView *create.CreateView, flags *pflag.FlagSet) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateLabelView(createView)

	if createView.Scope == "p" && !flags.Changed("project") {
		projectID, err := prompt.GetProjectIDFromUser()
		if err != nil {
			return fmt.Errorf("failed to get project id: %v", err)
		}

		createView.ProjectID = projectID
	}
	return api.CreateLabel(*createView)
}
