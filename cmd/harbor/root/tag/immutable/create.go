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
package immutable

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/immutable/create"
	"github.com/spf13/cobra"
)

func CreateImmutableCommand() *cobra.Command {
	var opts create.CreateView

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "create immutable tag rule",
		Long:    "create immutable tag rule to the project in harbor",
		Args:    cobra.MaximumNArgs(1),
		Example: "harbor tag immutable create",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName string
			createView := &create.CreateView{
				ScopeSelectors: create.ImmutableSelector{
					Decoration: opts.ScopeSelectors.Decoration,
					Pattern:    opts.ScopeSelectors.Pattern,
				},
				TagSelectors: create.ImmutableSelector{
					Decoration: opts.TagSelectors.Decoration,
					Pattern:    opts.TagSelectors.Pattern,
				},
			}
			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			err = createImmutableView(createView, projectName)
			if err != nil {
				return fmt.Errorf("failed to create immutable tag rule: %v", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.ScopeSelectors.Decoration, "repo-decoration", "", "", "repository which either apply or exclude from the rule")
	flags.StringVarP(&opts.ScopeSelectors.Pattern, "repo-list", "", "", "list of repository to which to either apply or exclude from the rule")
	flags.StringVarP(&opts.TagSelectors.Decoration, "tag-decoration", "", "", "tags which either apply or exclude from the rule")
	flags.StringVarP(&opts.TagSelectors.Pattern, "tag-list", "", "", "list of tags to which to either apply or exclude from the rule")

	return cmd
}

func createImmutableView(createView *create.CreateView, projectName string) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateImmutableView(createView)
	return api.CreateImmutable(*createView, projectName)
}
