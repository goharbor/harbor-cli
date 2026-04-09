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
	"errors"
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/immutable/create"
	"github.com/spf13/cobra"
)

func UpdateImmutableCommand() *cobra.Command {
	var opts create.CreateView
	var immutableID int64

	cmd := &cobra.Command{
		Use:   "update [PROJECT_NAME]",
		Short: "update immutable tag rule",
		Long:  "update immutable tag rule for a project in harbor",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName string

			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			if immutableID == 0 {
				immutableID = prompt.GetImmutableTagRule(projectName)
			}

			existingRule, err := getImmutableRuleByID(projectName, immutableID)
			if err != nil {
				return fmt.Errorf("failed to load existing immutable tag rule: %v", err)
			}

			updateView := buildUpdateViewFromRule(existingRule)
			applyFlagOverrides(updateView, opts)

			if !hasAnyUpdateFlag(opts) {
				create.CreateImmutableView(updateView)
			}

			applyViewToRule(existingRule, updateView)

			if existingRule.ID == 0 {
				existingRule.ID = immutableID
			}

			err = api.UpdateImmutable(existingRule, projectName, immutableID)
			if err != nil {
				return fmt.Errorf("failed to update immutable tag rule: %v", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&immutableID, "immutable-id", "", 0, "immutable rule ID to update")
	flags.StringVarP(&opts.ScopeSelectors.Decoration, "repo-decoration", "", "", "repository which either apply or exclude from the rule")
	flags.StringVarP(&opts.ScopeSelectors.Pattern, "repo-list", "", "", "list of repository to which to either apply or exclude from the rule")
	flags.StringVarP(&opts.TagSelectors.Decoration, "tag-decoration", "", "", "tags which either apply or exclude from the rule")
	flags.StringVarP(&opts.TagSelectors.Pattern, "tag-list", "", "", "list of tags to which to either apply or exclude from the rule")

	return cmd
}

func getImmutableRuleByID(projectName string, immutableID int64) (*models.ImmutableRule, error) {
	resp, err := api.ListImmutable(projectName)
	if err != nil {
		return nil, err
	}

	for _, rule := range resp.Payload {
		if rule != nil && rule.ID == immutableID {
			return rule, nil
		}
	}

	return nil, errors.New("immutable rule not found")
}

func buildUpdateViewFromRule(rule *models.ImmutableRule) *create.CreateView {
	view := &create.CreateView{}
	if rule == nil {
		return view
	}

	if len(rule.TagSelectors) > 0 && rule.TagSelectors[0] != nil {
		view.TagSelectors.Decoration = rule.TagSelectors[0].Decoration
		view.TagSelectors.Pattern = rule.TagSelectors[0].Pattern
	}

	if repoScopes, ok := rule.ScopeSelectors["repository"]; ok && len(repoScopes) > 0 {
		view.ScopeSelectors.Decoration = repoScopes[0].Decoration
		view.ScopeSelectors.Pattern = repoScopes[0].Pattern
	}

	return view
}

func applyFlagOverrides(target *create.CreateView, opts create.CreateView) {
	if target == nil {
		return
	}

	if opts.ScopeSelectors.Decoration != "" {
		target.ScopeSelectors.Decoration = opts.ScopeSelectors.Decoration
	}
	if opts.ScopeSelectors.Pattern != "" {
		target.ScopeSelectors.Pattern = opts.ScopeSelectors.Pattern
	}
	if opts.TagSelectors.Decoration != "" {
		target.TagSelectors.Decoration = opts.TagSelectors.Decoration
	}
	if opts.TagSelectors.Pattern != "" {
		target.TagSelectors.Pattern = opts.TagSelectors.Pattern
	}
}

func applyViewToRule(rule *models.ImmutableRule, view *create.CreateView) {
	if rule == nil || view == nil {
		return
	}

	rule.TagSelectors = []*models.ImmutableSelector{{
		Decoration: view.TagSelectors.Decoration,
		Pattern:    view.TagSelectors.Pattern,
	}}

	rule.ScopeSelectors = map[string][]models.ImmutableSelector{
		"repository": {
			{
				Decoration: view.ScopeSelectors.Decoration,
				Pattern:    view.ScopeSelectors.Pattern,
			},
		},
	}
}

func hasAnyUpdateFlag(opts create.CreateView) bool {
	return opts.ScopeSelectors.Decoration != "" ||
		opts.ScopeSelectors.Pattern != "" ||
		opts.TagSelectors.Decoration != "" ||
		opts.TagSelectors.Pattern != ""
}
