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
package retention

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/retention/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateRetentionCommand() *cobra.Command {
	var opts create.CreateView
	var projectName string
	var projectID int
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a tag retention rule in a project",
		Long: `Create a tag retention rule for a project in Harbor to manage the lifecycle of image tags.

Tag retention rules help users automatically retain or delete specific tags based on 
defined criteria, reducing storage usage and improving repository maintenance.

A user can create up to 15 tag retention rules per project.`,
		Example: `  # Retain tags matching 'release-*' at the project level
  harbor tag retention create --level project --action retain --taglist release-*

  # Delete untagged images at the repository level
  harbor retention create --level repository --action delete --tagdecoration untagged`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if projectID != -1 && projectName != "" {
				return fmt.Errorf("Cannot specify both --project-id and --project-name flags")
			}

			if projectID == -1 && projectName == "" {
				projectName = prompt.GetProjectNameFromUser()
			}

			projectIDorName := ""
			isName := true
			if projectID != -1 {
				projectIDorName = strconv.Itoa(projectID)
				isName = false
			} else {
				projectIDorName = projectName
			}

			createView := &create.CreateView{
				ScopeSelectors: create.RetentionSelector{
					Decoration: opts.ScopeSelectors.Decoration,
					Pattern:    opts.ScopeSelectors.Pattern,
				},
				TagSelectors: create.RetentionSelector{
					Decoration: opts.TagSelectors.Decoration,
					Pattern:    opts.TagSelectors.Pattern,
					Extras:     opts.TagSelectors.Extras,
				},
				Scope: create.RetentionPolicyScope{
					Level: opts.Scope.Level,
					Ref:   opts.Scope.Ref,
				},
				Template:  opts.Template,
				Params:    opts.Params,
				Action:    opts.Action,
				Algorithm: opts.Algorithm,
			}

			projectId, err := prompt.GetProjectIDFromUser()
			if err != nil {
				return err
			}

			err = createRetentionView(createView, int32(projectId))
			if err != nil {
				log.Errorf("Failed to create tag retention rule: %v", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.ScopeSelectors.Decoration, "repodecoration", "", "", "Apply or exclude repositories from the rule")
	flags.StringVarP(&opts.ScopeSelectors.Pattern, "repolist", "", "", "Comma-separated list of repositories to apply/exclude")
	flags.StringVarP(&opts.TagSelectors.Decoration, "tagdecoration", "", "", "Apply or exclude specific tags from the rule")
	flags.StringVarP(&opts.TagSelectors.Pattern, "taglist", "", "", "Comma-separated list of tags to apply/exclude")
	flags.StringVarP(&opts.Scope.Level, "level", "", "project", "Scope of the retention policy: 'project' or 'repository'")
	flags.StringVarP(&opts.Action, "action", "", "retain", "Action to perform: 'retain' or 'delete'")
	flags.StringVarP(&opts.Algorithm, "algorithm", "", "or", "Rule combination method: 'or' or 'and'")
	flags.StringVarP(&projectName, "project-name", "p", "", "Project name")
	flags.IntVarP(&projectID, "project-id", "i", -1, "Project ID")

	return cmd
}

func createRetentionView(createView *create.CreateView, projectIDorName string, isName bool) error {
	if createView == nil {
		createView = &create.CreateView{}
	}

	create.CreateRetentionView(createView)
	return api.CreateRetention(*createView, projectIDorName, isName)
}
