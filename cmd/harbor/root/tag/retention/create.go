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
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	rcreate "github.com/goharbor/harbor-cli/pkg/views/retention/create"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	var policyFile string
	var projectName string
	var dryRun bool
	var opts rcreate.CreateView

	cmd := &cobra.Command{
		Use:   "create [PROJECT_NAME]",
		Short: "create retention policy",
		Long:  "create a retention policy for a project",
		Args:  cobra.MaximumNArgs(1),
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

			var policy *models.RetentionPolicy
			var err error
			if policyFile != "" {
				policy, err = loadRetentionPolicyFromFile(policyFile)
				if err != nil {
					return err
				}
			} else {
				view := opts
				if !hasAnyCreateFlag(opts) {
					rcreate.CreateRetentionView(&view)
				}
				policy = buildRetentionPolicyFromView(view)
			}

			projectID, err := api.GetProjectIDFromName(projectName)
			if err != nil {
				return fmt.Errorf("failed to resolve project ID for %q: %v", projectName, utils.ParseHarborErrorMsg(err))
			}

			if policy.Scope == nil {
				policy.Scope = &models.RetentionPolicyScope{}
			}
			policy.Scope.Level = "project"
			policy.Scope.Ref = projectID

			if dryRun {
				if policyFile != "" {
					fmt.Printf("Retention policy file %q validated for project %q (dry-run).\n", policyFile, projectName)
				} else {
					fmt.Printf("Retention policy validated for project %q (dry-run).\n", projectName)
				}
				return nil
			}

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

	flags := cmd.Flags()
	flags.StringVarP(&policyFile, "file", "f", "", "retention policy file in JSON format (optional in interactive mode)")
	flags.StringVarP(&projectName, "project", "", "", "project name")
	flags.BoolVarP(&dryRun, "dry-run", "", false, "validate policy file and project scope without creating")
	flags.StringVarP(&opts.ScopeSelectors.Decoration, "repo-decoration", "", "", "repository selector decoration: repoMatches or repoExcludes")
	flags.StringVarP(&opts.ScopeSelectors.Pattern, "repo-list", "", "", "repository selector pattern, for example **")
	flags.StringVarP(&opts.TagSelectors.Decoration, "tag-decoration", "", "", "tag selector decoration: matches or excludes")
	flags.StringVarP(&opts.TagSelectors.Pattern, "tag-list", "", "", "tag selector pattern, for example **")
	flags.Int64VarP(&opts.KeepLatestPushed, "keep-latest", "", 0, "number of most recently pushed artifacts to retain")
	flags.StringVarP(&opts.Cron, "cron", "", "", "schedule cron expression for retention policy")

	return cmd
}

func loadRetentionPolicyFromFile(filePath string) (*models.RetentionPolicy, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file: %v", err)
	}

	policy := &models.RetentionPolicy{}
	if err := json.Unmarshal(data, policy); err != nil {
		return nil, fmt.Errorf("failed to parse policy JSON: %v", err)
	}

	if len(policy.Rules) == 0 {
		return nil, fmt.Errorf("policy must contain at least one rule")
	}

	return policy, nil
}

func hasAnyCreateFlag(opts rcreate.CreateView) bool {
	return opts.ScopeSelectors.Decoration != "" ||
		opts.ScopeSelectors.Pattern != "" ||
		opts.TagSelectors.Decoration != "" ||
		opts.TagSelectors.Pattern != "" ||
		opts.KeepLatestPushed > 0 ||
		opts.Cron != ""
}

func buildRetentionPolicyFromView(view rcreate.CreateView) *models.RetentionPolicy {
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
		Algorithm: "or",
		Rules: []*models.RetentionRule{
			{
				Action:   "retain",
				Template: "latestPushedK",
				Params: map[string]interface{}{
					"latestPushedK": view.KeepLatestPushed,
				},
				ScopeSelectors: map[string][]models.RetentionSelector{
					"repository": {
						{
							Kind:       "doublestar",
							Decoration: view.ScopeSelectors.Decoration,
							Pattern:    view.ScopeSelectors.Pattern,
						},
					},
				},
				TagSelectors: []*models.RetentionSelector{
					{
						Kind:       "doublestar",
						Decoration: view.TagSelectors.Decoration,
						Pattern:    view.TagSelectors.Pattern,
					},
				},
			},
		},
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
