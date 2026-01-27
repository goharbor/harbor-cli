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

	"github.com/goharbor/harbor-cli/pkg/views/replication/policies/create"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"

)

type PolicyConfig struct {
	Name              string               `yaml:"name" json:"name"`
	Description       string               `yaml:"description" json:"description"`
	ReplicationMode   string               `yaml:"replication_mode,omitempty" json:"replication_mode,omitempty"`
	Filter            []*ReplicationFilter `yaml:"replication_filter,omitempty" json:"replication_filter,omitempty"`
	TargetRegistry    string               `yaml:"target_registry,omitempty" json:"target_registry,omitempty"`
	TriggerMode       string               `yaml:"trigger_mode,omitempty" json:"trigger_mode,omitempty"`
	BandWidthLimit    string               `yaml:"bandwidth_limit,omitempty" json:"bandwidth_limit,omitempty"`
	CronString        string               `yaml:"cron_string,omitempty" json:"cron_string,omitempty"`
	Override          bool                 `yaml:"override,omitempty" json:"override,omitempty"`
	ReplicateDeletion bool                 `yaml:"replicate_deletion,omitempty" json:"replicate_deletion,omitempty"`
	CopyByChunk       bool                 `yaml:"copy_by_chunk,omitempty" json:"copy_by_chunk,omitempty"`
	Enabled           bool                 `yaml:"enabled,omitempty" json:"enabled,omitempty"`
}

type ReplicationFilter struct {
	Type       string `yaml:"type,omitempty" json:"type,omitempty"`
	Decoration string `yaml:"decoration,omitempty" json:"decoration,omitempty"`
	Value      string `yaml:"value,omitempty" json:"value,omitempty"`
}

func LoadConfigFromFile(filename string) (*create.CreateView, error) {
	var opts *create.CreateView
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

func LoadConfigFromYAMLorJSON(filename string, fileType string) (*create.CreateView, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	log.Debug("Replication policy config file read successfully")


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
	log.Debug("Replication policy configuration validated successfully")


	opts := &create.CreateView{
		Name:              config.Name,
		Description:       config.Description,
		ReplicationMode:   normalizeReplicationMode(config.ReplicationMode),
		TriggerType:       normalizeTriggerMode(config.TriggerMode),
		TargetRegistry:    config.TargetRegistry,
		CronString:        config.CronString,
		Override:          config.Override,
		CopyByChunk:       config.CopyByChunk,
		ReplicateDeletion: config.ReplicateDeletion,
		Speed:             config.BandWidthLimit,
		Enabled:           config.Enabled,
	}

	if err := processFilters(&config, opts); err != nil {
		return nil, fmt.Errorf("failed to process filters: %v", err)
	}

	return opts, nil
}

func validateConfig(config *PolicyConfig) error {
	if config.Name == "" {
		return fmt.Errorf("name is required")
	}

	if config.ReplicationMode != "" {
		mode := strings.ToLower(config.ReplicationMode)
		if mode != "push" && mode != "pull" {
			return fmt.Errorf("replication_mode must be 'push' or 'pull', got: %s", config.ReplicationMode)
		}
	}

	if config.TriggerMode != "" {
		mode := strings.ToLower(config.TriggerMode)
		validTriggers := []string{"manual", "scheduled", "event_based"}
		isValid := false
		for _, valid := range validTriggers {
			if mode == valid {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("trigger_mode must be one of [manual, scheduled, event_based], got: %s", config.TriggerMode)
		}
	}

	for i, filter := range config.Filter {
		if err := validateFilter(filter, i); err != nil {
			return err
		}
	}

	return nil
}

func validateFilter(filter *ReplicationFilter, index int) error {
	if filter.Type == "" {
		return fmt.Errorf("filter[%d]: type is required", index)
	}

	validTypes := []string{"resource", "name", "tag", "label"}
	isValidType := false
	for _, validType := range validTypes {
		if filter.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("filter[%d]: type must be one of [resource, name, tag, label], got: %s", index, filter.Type)
	}

	if filter.Type == "resource" {
		if filter.Value != "" {
			validResources := []string{"image", "artifact"}
			isValidResource := false
			for _, validResource := range validResources {
				if filter.Value == validResource {
					isValidResource = true
					break
				}
			}
			if !isValidResource {
				return fmt.Errorf("filter[%d]: resource value must be 'image' or 'chart', got: %s", index, filter.Value)
			}
		}
	}

	if filter.Decoration != "" {
		if filter.Type != "tag" && filter.Type != "label" {
			return fmt.Errorf("filter[%d]: decoration is only supported for 'tag' and 'label' filters, got type: %s", index, filter.Type)
		}

		validDecorations := []string{"matches", "excludes"}
		isValidDecoration := false
		for _, validDecoration := range validDecorations {
			if filter.Decoration == validDecoration {
				isValidDecoration = true
				break
			}
		}
		if !isValidDecoration {
			return fmt.Errorf("filter[%d]: decoration must be 'matches' or 'excludes', got: %s", index, filter.Decoration)
		}
	}

	return nil
}

func normalizeReplicationMode(mode string) string {
	switch strings.ToLower(mode) {
	case "push":
		return "Push"
	case "pull":
		return "Pull"
	default:
		return mode
	}
}

func normalizeTriggerMode(mode string) string {
	switch strings.ToLower(mode) {
	case "manual":
		return "manual"
	case "scheduled":
		return "scheduled"
	case "event_based":
		return "event_based"
	default:
		return mode
	}
}

func processFilters(config *PolicyConfig, opts *create.CreateView) error {
	for _, filter := range config.Filter {
		switch filter.Type {
		case "resource":
			if filter.Value == "" {
				opts.ResourceFilter = "All"
			} else {
				opts.ResourceFilter = filter.Value
			}
		case "name":
			opts.NameFilter = filter.Value
		case "tag":
			if filter.Decoration != "" {
				opts.TagFilter = filter.Decoration
			} else {
				opts.TagFilter = "matches"
			}
			opts.TagPattern = filter.Value
		case "label":
			if filter.Decoration != "" {
				opts.LabelFilter = filter.Decoration
			} else {
				opts.LabelFilter = "matches"
			}
			opts.LabelPattern = filter.Value
		}
	}

	return nil
}
