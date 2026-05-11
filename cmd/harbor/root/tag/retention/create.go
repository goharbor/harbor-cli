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
	"errors"
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	rcreate "github.com/goharbor/harbor-cli/pkg/views/retention/create"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	var projectName string

	cmd := &cobra.Command{
		Use:   "create [PROJECT_NAME]",
		Short: "create retention policy",
		Long:  "create a retention policy for a project",
		Example: `
# Create a retention policy for a specific project
harbor tag retention create my-project

# Create a retention policy interactively
harbor tag retention create`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				projectName = args[0]
			}

			if projectName == "" {
				var err error
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			view := rcreate.CreateView{}
			rcreate.CreateRetentionView(&view)
			policy := buildRetentionPolicyFromView(view)

			projectID, err := api.GetProjectIDFromName(projectName)
			if err != nil {
				return fmt.Errorf("failed to resolve project ID for %q: %v", projectName, utils.ParseHarborErrorMsg(err))
			}

			if policy.Scope == nil {
				policy.Scope = &models.RetentionPolicyScope{}
			}
			policy.Scope.Level = "project"
			policy.Scope.Ref = projectID

			existingRetentionID, err := api.GetRetentionIDByProjectName(projectName)
			if err == nil {
				existingPolicy, err := api.GetRetentionPolicy(existingRetentionID)
				if err != nil {
					return fmt.Errorf("failed to get existing retention policy %d: %v", existingRetentionID, utils.ParseHarborErrorMsg(err))
				}

				mergedPolicy := mergeRuleIntoExistingPolicy(existingPolicy, policy)
				sanitizeRetentionPolicyForUpdate(mergedPolicy)

				if err := api.UpdateRetention(existingRetentionID, mergedPolicy); err != nil {
					return fmt.Errorf("failed to append rule to retention policy %d: %v", existingRetentionID, utils.ParseHarborErrorMsg(err))
				}

				fmt.Printf("Retention rule added successfully to existing policy %d\n", existingRetentionID)
				return nil
			}

			if !errors.Is(err, api.ErrNoRetentionPolicy) {
				return fmt.Errorf("failed to check existing retention policy for project %q: %v", projectName, utils.ParseHarborErrorMsg(err))
			}

			location, err := api.CreateRetention(policy)
			if err != nil {
				return fmt.Errorf("failed to create retention policy: %v", utils.ParseHarborErrorMsg(err))
			}

			if location != "" {
				fmt.Printf("Retention policy created successfully: %s\n", location)
				return nil
			}

			fmt.Println("Retention policy created successfully")
			return nil
		},
	}

	return cmd
}

func buildRetentionPolicyFromView(view rcreate.CreateView) *models.RetentionPolicy {
	// NOTE: The retention policy structure is defined by the Harbor API (POST /api/v2.0/retentions).

	// Key constraints (as enforced by Harbor backend):
	// - algorithm: "or" (current UI does not support "and")
	// - action: "retain" (artifacts matching selectors will be kept)
	// - template: "latestPushedK" (keeps K most recently pushed artifacts)
	// - kind: "doublestar" (supports glob patterns: *, **, etc.)
	// - scope_selectors decorations: "repoMatches" | "repoExcludes"
	// - tag_selectors decorations: "matches" | "excludes"
	// - trigger.kind: "Schedule" (with cron expression in settings)

	if view.ScopeSelectors.Decoration == "" {
		view.ScopeSelectors.Decoration = "repoMatches"
	}
	if view.ScopeSelectors.Pattern == "" {
		view.ScopeSelectors.Pattern = "**"
	}
	if view.TagSelectors.Decoration == "" {
		view.TagSelectors.Decoration = "matches"
	}
	if view.TagSelectors.Pattern == "" {
		view.TagSelectors.Pattern = "**"
	}
	if view.KeepLatestPushed <= 0 {
		view.KeepLatestPushed = 10
	}
	if view.Cron == "" {
		view.Cron = "0 0 0 * * *"
	}

	return &models.RetentionPolicy{
		// NOTE: Algorithm is hardcoded to "or" because the current UI does not support selecting "and".
		// "or" is currently the only practical/useful decision for retention policies.
		// If "and" algorithm support is needed in the future, both the UI and this hardcoded value
		// would need to be updated to allow users to choose the algorithm.
		Algorithm: "or",
		Rules: []*models.RetentionRule{
			{
				// Action "retain" means artifacts matching these selectors will be kept
				Action: "retain",
				// Template "latestPushedK" retains the K most recently pushed artifacts
				Template: "latestPushedK",
				Params: map[string]interface{}{
					"latestPushedK": view.KeepLatestPushed,
				},
				// Repository selector using doublestar pattern matching
				// kind: "doublestar" supports glob patterns (*, **, domain/**, etc.)
				// decoration: "repoMatches" includes repos matching pattern, "repoExcludes" excludes them
				ScopeSelectors: map[string][]models.RetentionSelector{
					"repository": {
						{
							Kind:       "doublestar",
							Decoration: view.ScopeSelectors.Decoration,
							Pattern:    view.ScopeSelectors.Pattern,
						},
					},
				},
				// Tag selector with same doublestar matching
				// decoration: "matches" applies rule to tags matching pattern, "excludes" excludes them
				TagSelectors: []*models.RetentionSelector{
					{
						Kind:       "doublestar",
						Decoration: view.TagSelectors.Decoration,
						Pattern:    view.TagSelectors.Pattern,
					},
				},
			},
		},
		// Trigger with Schedule kind executes retention policy on a cron schedule
		// cron format: second minute hour day-of-month month day-of-week (6 fields)
		Trigger: &models.RetentionRuleTrigger{
			Kind:       "Schedule",
			Settings:   map[string]interface{}{"cron": view.Cron},
			References: map[string]interface{}{},
		},
	}
}

func mergeRuleIntoExistingPolicy(existing *models.RetentionPolicy, incoming *models.RetentionPolicy) *models.RetentionPolicy {
	if existing == nil {
		return incoming
	}

	if incoming == nil || len(incoming.Rules) == 0 {
		return existing
	}

	existing.Rules = append(existing.Rules, incoming.Rules...)
	if existing.Algorithm == "" {
		existing.Algorithm = "or"
	}

	return existing
}

func sanitizeRetentionPolicyForUpdate(policy *models.RetentionPolicy) {
	if policy == nil || policy.Trigger == nil || policy.Trigger.Settings == nil {
		return
	}

	if settings, ok := policy.Trigger.Settings.(map[string]interface{}); ok {
		delete(settings, "next_scheduled_time")
	}
}
