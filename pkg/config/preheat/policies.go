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
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	policycreate "github.com/goharbor/harbor-cli/pkg/views/preheat/policy/create"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type PolicyConfig struct {
	ProjectName  string          `yaml:"project_name" json:"project_name"`
	Name         string          `yaml:"name" json:"name"`
	Description  string          `yaml:"description,omitempty" json:"description,omitempty"`
	Enabled      bool            `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Filters      []*PolicyFilter `yaml:"filters" json:"filters"`
	ProviderName string          `yaml:"provider_name" json:"provider_name"`
	Trigger      *PolicyTrigger  `yaml:"trigger,omitempty" json:"trigger,omitempty"`
}

type PolicyFilter struct {
	Type  string `yaml:"type" json:"type"`
	Value string `yaml:"value" json:"value"`
}

type PolicyTriggerSetting struct {
	SchedulePreset string `yaml:"schedule_preset,omitempty" json:"schedule_preset,omitempty"`
	Cron           string `yaml:"cron,omitempty" json:"cron,omitempty"`
}

type PolicyTrigger struct {
	Type           string                `yaml:"type" json:"type"`
	TriggerSetting *PolicyTriggerSetting `yaml:"trigger_setting,omitempty" json:"trigger_setting,omitempty"`
}

func LoadConfigFromFile(filename string) (*policycreate.CreateView, error) {
	var opts *policycreate.CreateView
	var err error

	ext := filepath.Ext(filename)
	if ext == "" {
		return nil, fmt.Errorf("file must have an extension (.yaml, .yml, or .json)")
	}

	fileType := ext[1:]
	if fileType == "yml" {
		fileType = "yaml"
	}

	opts, err = LoadConfigFromYAMLorJSON(filename, fileType)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}
	return opts, nil
}

func LoadConfigFromYAMLorJSON(filename string, fileType string) (*policycreate.CreateView, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	log.Debug("Preheat policy config file read successfully")

	var config PolicyConfig
	switch fileType {
	case "yaml", "yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %v", err)
		}
		log.Debugf("Parsed %s configuration successfully", fileType)

	case "json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", err)
		}
		log.Debugf("Parsed %s configuration successfully", fileType)
	default:
		return nil, fmt.Errorf("unsupported file type: %s, expected 'yaml' or 'json'", fileType)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %v", err)
	}
	log.Debug("Preheat policy configuration validated successfully")

	triggerType := "manual"
	cronString := ""
	if config.Trigger != nil {
		triggerType = normalizeTriggerMode(config.Trigger.Type)
		if triggerType == "scheduled" {
			var err error
			cronString, err = resolveTriggerCron(config.Trigger)
			if err != nil {
				return nil, err
			}
		}
	}

	opts := &policycreate.CreateView{
		ProjectName:  config.ProjectName,
		Name:         config.Name,
		Description:  config.Description,
		ProviderName: config.ProviderName,
		TriggerType:  triggerType,
		CronString:   cronString,
		Enabled:      config.Enabled,
	}

	for _, filter := range config.Filters {
		switch filter.Type {
		case "repository":
			opts.RepositoryFilter = filter.Value
		case "tag":
			opts.TagFilter = filter.Value
		case "label":
			opts.LabelFilter = filter.Value
		}
	}

	return opts, nil
}

func validateConfig(config *PolicyConfig) error {
	if config.ProjectName == "" {
		return fmt.Errorf("project_name is required")
	}

	if config.Name == "" {
		return fmt.Errorf("name is required")
	}

	if config.ProviderName == "" {
		return fmt.Errorf("provider_name is required")
	}

	if err := validatePolicyFilters(config.Filters); err != nil {
		return err
	}

	if config.Trigger == nil {
		return nil
	}

	if strings.TrimSpace(config.Trigger.Type) == "" {
		return fmt.Errorf("trigger.type is required")
	}

	triggerType := strings.ToLower(config.Trigger.Type)
	switch triggerType {
	case "manual", "scheduled", "event_based":
	default:
		return fmt.Errorf("trigger.type must be one of [manual, scheduled, event_based], got: %s", config.Trigger.Type)
	}

	if triggerType != "scheduled" {
		setting := config.Trigger.TriggerSetting
		if setting != nil && (strings.TrimSpace(setting.SchedulePreset) != "" || strings.TrimSpace(setting.Cron) != "") {
			return fmt.Errorf("trigger.trigger_setting is only supported for scheduled trigger")
		}
		return nil
	}

	_, err := resolveTriggerCron(config.Trigger)
	return err
}

func resolveTriggerCron(trigger *PolicyTrigger) (string, error) {
	if trigger == nil || trigger.TriggerSetting == nil {
		return "", nil
	}

	setting := trigger.TriggerSetting
	preset := strings.ToLower(strings.TrimSpace(setting.SchedulePreset))
	switch preset {
	case "":
		return strings.TrimSpace(setting.Cron), nil
	case "none":
		return "", nil
	case "hourly", "daily", "weekly", "custom":
		cron := policycreate.ResolveSchedulePreset(preset, setting.Cron)
		if preset == "custom" && cron == "" {
			return "", fmt.Errorf("trigger.trigger_setting.cron is required for custom schedule")
		}
		return cron, nil
	default:
		return "", fmt.Errorf("trigger.trigger_setting.schedule_preset must be one of [none, hourly, daily, weekly, custom], got: %s", setting.SchedulePreset)
	}
}

func validatePolicyFilters(filters []*PolicyFilter) error {
	if len(filters) < 2 || len(filters) > 3 {
		return fmt.Errorf("filters must include repository and tag filters, with label optional")
	}

	hasRepositoryFilter := false
	hasTagFilter := false
	for i, filter := range filters {
		if filter == nil {
			return fmt.Errorf("filters[%d] is required", i)
		}

		switch filter.Type {
		case "repository":
			hasRepositoryFilter = true
		case "tag":
			hasTagFilter = true
		case "label":
		default:
			return fmt.Errorf("filters[%d].type must be one of [repository, tag, label], got: %s", i, filter.Type)
		}

		if strings.TrimSpace(filter.Value) == "" {
			return fmt.Errorf("filters[%d].value is required", i)
		}
	}

	if !hasRepositoryFilter {
		return fmt.Errorf("filters must include a repository filter")
	}
	if !hasTagFilter {
		return fmt.Errorf("filters must include a tag filter")
	}

	return nil
}

func normalizeTriggerMode(mode string) string {
	switch strings.ToLower(mode) {
	case "scheduled":
		return "scheduled"
	case "event_based":
		return "event_based"
	default:
		return "manual"
	}
}
