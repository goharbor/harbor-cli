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
package policy

import (
	"encoding/json"
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/preheat/policy/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdatePolicyCommand() *cobra.Command {
	var isID bool

	cmd := &cobra.Command{
		Use:     "update [NAME|ID] [POLICY_NAME]",
		Short:   "Update a preheat policy",
		Long:    "Update an existing P2P preheat policy under a project",
		Example: `  harbor-cli project preheat policy update [NAME|ID] [POLICY_NAME]`,
		Args:    cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName, policyName string

			if isID && len(args) == 0 {
				return fmt.Errorf("project ID must be provided when using --id")
			}

			if len(args) >= 1 {
				log.Debugf("Project name provided: %s", args[0])
				projectName = args[0]
			} else {
				log.Debug("No project name provided, prompting user")
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			if isID {
				project, err := api.GetProject(projectName, true)
				if err != nil {
					return fmt.Errorf("failed to get project: %v", utils.ParseHarborErrorMsg(err))
				}
				projectName = project.Payload.Name
			}

			if len(args) >= 2 {
				log.Debugf("Policy name provided: %s", args[1])
				policyName = args[1]
			} else {
				log.Debug("No policy name provided, prompting user")
				policyName, err = prompt.GetPreheatPolicyNameFromUser(projectName)
				if err != nil {
					return fmt.Errorf("failed to get policy name: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			log.Debug("Fetching preheat policy...")
			existingPolicy, err := api.GetPreheatPolicy(projectName, policyName)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "404" {
					return fmt.Errorf("preheat policy %s not found in project %s", policyName, projectName)
				}
				return fmt.Errorf("failed to get preheat policy: %v", utils.ParseHarborErrorMsg(err))
			}

			log.Debug("Fetching available providers...")
			providers, err := api.ListProvidersUnderProject(projectName)
			if err != nil {
				return fmt.Errorf("failed to list providers: %v", utils.ParseHarborErrorMsg(err))
			}

			if len(providers) == 0 {
				return fmt.Errorf("no P2P provider instances available for project '%s'. Please create a provider instance first", projectName)
			}

			opts := policyToCreateView(existingPolicy.Payload, providers)
			create.CreatePreheatPolicyView(opts, providers)

			providerID, err := resolveProviderID(providers, opts.ProviderName, projectName)
			if err != nil {
				return err
			}

			policy, err := ConvertToPolicy(opts, providerID)
			if err != nil {
				return err
			}
			policy.ID = existingPolicy.Payload.ID
			policy.ProjectID = existingPolicy.Payload.ProjectID
			policy.CreationTime = existingPolicy.Payload.CreationTime
			policy.ExtraAttrs = existingPolicy.Payload.ExtraAttrs

			log.Debug("Updating preheat policy...")
			_, err = api.UpdatePreheatPolicy(projectName, policyName, policy)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "409" {
					return fmt.Errorf("preheat policy '%s' already exists in project '%s'", opts.Name, projectName)
				}
				return fmt.Errorf("failed to update preheat policy: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Printf("Preheat policy '%s' updated successfully in project '%s'\n", opts.Name, projectName)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&isID, "id", false, "Use project id instead of name")

	return cmd
}

func policyToCreateView(policy *models.PreheatPolicy, providers []*models.ProviderUnderProject) *create.CreateView {
	view := &create.CreateView{
		Name:         policy.Name,
		Description:  policy.Description,
		Enabled:      policy.Enabled,
		ProviderName: policy.ProviderName,
		TriggerType:  "manual",
	}

	for _, provider := range providers {
		if provider.ID == policy.ProviderID {
			view.ProviderName = provider.Provider
			break
		}
	}

	var filters []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	_ = json.Unmarshal([]byte(policy.Filters), &filters)
	for _, filter := range filters {
		switch filter.Type {
		case "repository":
			view.RepositoryFilter = filter.Value
		case "tag":
			view.TagFilter = filter.Value
		case "label":
			view.LabelFilter = filter.Value
		}
	}

	type triggerSetting struct {
		Cron string `json:"cron"`
	}
	var trigger struct {
		Type            string          `json:"type"`
		TriggerSetting  *triggerSetting `json:"trigger_setting"`
		TriggerSettings *triggerSetting `json:"trigger_settings"`
	}
	_ = json.Unmarshal([]byte(policy.Trigger), &trigger)
	if trigger.Type != "" {
		view.TriggerType = trigger.Type
	}
	if trigger.TriggerSetting != nil {
		view.CronString = trigger.TriggerSetting.Cron
	} else if trigger.TriggerSettings != nil {
		view.CronString = trigger.TriggerSettings.Cron
	}

	return view
}
