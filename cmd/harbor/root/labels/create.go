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
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/label/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateLabelCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create label",
		Long:    "create label in harbor",
		Example: "harbor label create",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			createView := &create.CreateView{
				Name:        opts.Name,
				Color:       opts.Color,
				Scope:       opts.Scope,
				Description: opts.Description,
			}
			if opts.Name != "" && opts.Scope != "" {
				err = api.CreateLabel(opts)
			} else {
				err = createLabelView(createView)
			}

			if err != nil {
				log.Errorf("failed to create label: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "Name of the label")
	flags.StringVarP(&opts.Color, "color", "", "#FFFFFF", "Color of the label.color is in hex value")
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "Scope of the label. eg- g(global), p(specific project)")
	flags.StringVarP(&opts.Description, "description", "d", "", "Description of the label")

	return cmd
}

func createLabelView(createView *create.CreateView) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateLabelView(createView)
	return api.CreateLabel(*createView)
}
