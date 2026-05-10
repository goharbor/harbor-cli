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
	config "github.com/goharbor/harbor-cli/pkg/config/preheat"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/preheat/policy/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreatePolicyCommand() *cobra.Command {
	var isID bool
	var configFile string

	cmd := &cobra.Command{
		Use:   "create [NAME|ID]",
		Short: "Create a preheat policy",
		Long:  "Create a new P2P preheat policy under a project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName string

			if len(args) > 0 {
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

			var opts *create.CreateView

			if configFile != "" {
				log.Debugf("Loading preheat policy configuration from file: %s", configFile)
				opts, err = config.LoadConfigFromFile(configFile)
				if err != nil {
					return fmt.Errorf("failed to load preheat policy configuration: %v", err)
				}
			} else {
				opts = &create.CreateView{
					Enabled: true,
				}
			}

			log.Debug("Fetching available providers...")
			providers, err := api.ListProvidersUnderProject(projectName)
			if err != nil {
				return fmt.Errorf("failed to list providers: %v", utils.ParseHarborErrorMsg(err))
			}

			if len(providers) == 0 {
				return fmt.Errorf("no P2P provider instances available for project '%s'. Please create a provider instance first", projectName)
			}

			if configFile == "" {
				create.CreatePreheatPolicyView(opts, providers)
			}

			providerID, err := resolveProviderID(providers, opts.ProviderName, projectName)
			if err != nil {
				return err
			}

			policy, err := ConvertToPolicy(opts, providerID)
			if err != nil {
				return err
			}

			log.Debug("Creating preheat policy...")
			response, err := api.CreatePreheatPolicy(projectName, policy)
			if err != nil {
				if utils.ParseHarborErrorCode(err) == "409" {
					return fmt.Errorf("preheat policy '%s' already exists in project '%s'", opts.Name, projectName)
				}
				return fmt.Errorf("failed to create preheat policy: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Println("Preheat policy created successfully with ID:", response.Location)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&isID, "id", false, "Use project id instead of name")
	flags.StringVarP(&configFile, "policy-config-file", "f", "", "YAML/JSON file with preheat policy configuration")

	return cmd
}

func resolveProviderID(providers []*models.ProviderUnderProject, providerName, projectName string) (int64, error) {
	for _, provider := range providers {
		if provider.Provider == providerName && provider.Enabled {
			return provider.ID, nil
		}
	}

	return 0, fmt.Errorf("provider '%s' not found or not enabled for project '%s'", providerName, projectName)
}

func ConvertToPolicy(view *create.CreateView, providerID int64) (*models.PreheatPolicy, error) {
	type filter struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	filters := []filter{
		{Type: "repository", Value: view.RepositoryFilter},
		{Type: "tag", Value: view.TagFilter},
	}
	if view.LabelFilter != "" {
		filters = append(filters, filter{Type: "label", Value: view.LabelFilter})
	}
	filtersJSON, err := json.Marshal(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filters: %v", err)
	}

	type triggerSetting struct {
		Cron string `json:"cron"`
	}
	type trigger struct {
		Type           string          `json:"type"`
		TriggerSetting *triggerSetting `json:"trigger_setting,omitempty"`
	}
	t := trigger{Type: view.TriggerType}
	if view.TriggerType == "scheduled" {
		t.TriggerSetting = &triggerSetting{Cron: view.CronString}
	}
	triggerJSON, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trigger: %v", err)
	}

	return &models.PreheatPolicy{
		Name:         view.Name,
		Description:  view.Description,
		ProviderID:   providerID,
		ProviderName: view.ProviderName,
		Filters:      string(filtersJSON),
		Trigger:      string(triggerJSON),
		Enabled:      view.Enabled,
	}, nil
}
